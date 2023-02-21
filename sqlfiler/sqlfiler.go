package sqlfiler

import (
	"bufio"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// InitializerImpl is the default implementation of the GrpcServer interface.
type SqlfilerImpl struct {
	DatabaseConnector driver.Connector
	isTrace           bool
	logger            messagelogger.MessageLoggerInterface
	LogLevel          logger.Level
	observers         subject.Subject
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (sqlfiler *SqlfilerImpl) getLogger() messagelogger.MessageLoggerInterface {
	if sqlfiler.logger == nil {
		sqlfiler.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return sqlfiler.logger
}

// Notify registered observers.
func (sqlfiler *SqlfilerImpl) notify(ctx context.Context, messageId int, err error, details map[string]string) {
	now := time.Now()
	details["subjectId"] = strconv.Itoa(ProductId)
	details["messageId"] = strconv.Itoa(messageId)
	details["messageTime"] = strconv.FormatInt(now.UnixNano(), 10)
	if err != nil {
		details["error"] = err.Error()
	}
	message, err := json.Marshal(details)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		sqlfiler.observers.NotifyObservers(ctx, string(message))
	}
}

// Trace method entry.
func (sqlfiler *SqlfilerImpl) traceEntry(errorNumber int, details ...interface{}) {
	sqlfiler.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (sqlfiler *SqlfilerImpl) traceExit(errorNumber int, details ...interface{}) {
	sqlfiler.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The ProcessFileName is a convenience method for calling method ProcessScanner using a filename.

Input
  - ctx: A context to control lifecycle.
  - filename: A fully qualified path to a file of SQL statements.
*/
func (sqlfiler *SqlfilerImpl) ProcessFileName(ctx context.Context, filename string) error {

	// Entry tasks.

	if sqlfiler.isTrace {
		sqlfiler.traceEntry(1, filename)
	}
	var err error = nil
	entryTime := time.Now()

	// Process file.

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	sqlfiler.ProcessScanner(ctx, bufio.NewScanner(file))

	// Exit tasks.

	if sqlfiler.observers != nil {
		go func() {
			details := map[string]string{
				"filename": filename,
			}
			sqlfiler.notify(ctx, 8001, err, details)
		}()
	}
	if sqlfiler.isTrace {
		defer sqlfiler.traceExit(2, filename, err, time.Since(entryTime))
	}
	return err
}

/*
The ProcessScanner does a database call for each line scanned.

Input
  - ctx: A context to control lifecycle.
  - scanner: SQL statements to be processed.
*/
func (sqlfiler *SqlfilerImpl) ProcessScanner(ctx context.Context, scanner *bufio.Scanner) {

	// Entry tasks.

	if sqlfiler.isTrace {
		sqlfiler.traceEntry(3)
	}
	var err error = nil
	entryTime := time.Now()

	// Open a database connection.

	database := sql.OpenDB(sqlfiler.DatabaseConnector)
	defer database.Close()

	// Process each scanned line.

	for scanner.Scan() {
		sqlText := scanner.Text()
		result, err := database.ExecContext(ctx, sqlText)
		if err != nil {
			sqlfiler.getLogger().Log(4001, result, err)
		}
		if sqlfiler.observers != nil {
			go func() {
				details := map[string]string{
					"SQL": sqlText,
				}
				sqlfiler.notify(ctx, 8002, err, details)
			}()
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Exit tasks.

	if sqlfiler.observers != nil {
		go func() {
			details := map[string]string{}
			sqlfiler.notify(ctx, 8003, err, details)
		}()
	}
	if sqlfiler.isTrace {
		defer sqlfiler.traceExit(4, err, time.Since(entryTime))
	}
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (sqlfiler *SqlfilerImpl) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	if sqlfiler.isTrace {
		sqlfiler.traceEntry(5, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	if sqlfiler.observers == nil {
		sqlfiler.observers = &subject.SubjectImpl{}
	}
	err := sqlfiler.observers.RegisterObserver(ctx, observer)
	if sqlfiler.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverId(ctx),
			}
			sqlfiler.notify(ctx, 8004, err, details)
		}()
	}
	if sqlfiler.isTrace {
		defer sqlfiler.traceExit(6, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (sqlfiler *SqlfilerImpl) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if sqlfiler.isTrace {
		sqlfiler.traceEntry(7, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	sqlfiler.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	sqlfiler.isTrace = (sqlfiler.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if sqlfiler.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logger.LevelToTextMap[logLevel],
			}
			sqlfiler.notify(ctx, 8005, err, details)
		}()
	}
	if sqlfiler.isTrace {
		defer sqlfiler.traceExit(8, logLevel, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (sqlfiler *SqlfilerImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	if sqlfiler.isTrace {
		sqlfiler.traceEntry(9, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	var err error = nil
	if sqlfiler.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		sqlfiler.notify(ctx, 8006, err, details)
	}
	err = sqlfiler.observers.UnregisterObserver(ctx, observer)
	if !sqlfiler.observers.HasObservers(ctx) {
		sqlfiler.observers = nil
	}
	if sqlfiler.isTrace {
		defer sqlfiler.traceExit(10, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}
