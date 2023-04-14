package senzingschema

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/senzing/go-common/engineconfigurationjsonparser"
	"github.com/senzing/go-databasing/connector"
	"github.com/senzing/go-databasing/sqlexecutor"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-observing/notifier"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// SenzingSchemaImpl is the default implementation of the SenzingSchema interface.
type SenzingSchemaImpl struct {
	isTrace                        bool
	logger                         logging.LoggingInterface
	logLevelName                   string
	observers                      subject.Subject
	SenzingEngineConfigurationJson string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var traceOptions []interface{} = []interface{}{
	&logging.OptionCallerSkip{Value: 5},
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (senzingSchema *SenzingSchemaImpl) getLogger() logging.LoggingInterface {
	var err error = nil
	if senzingSchema.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		senzingSchema.logger, err = logging.NewSenzingToolsLogger(ProductId, IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return senzingSchema.logger
}

// Log message.
func (senzingSchema *SenzingSchemaImpl) log(messageNumber int, details ...interface{}) {
	senzingSchema.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (senzingSchema *SenzingSchemaImpl) traceEntry(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	senzingSchema.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (senzingSchema *SenzingSchemaImpl) traceExit(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	senzingSchema.getLogger().Log(messageNumber, details...)
}

// --- Misc -------------------------------------------------------------------

// Given a database URL, detemine the correct SQL file and send the statements to the database.
func (senzingSchema *SenzingSchemaImpl) processDatabase(ctx context.Context, resourcePath string, databaseUrl string) error {
	var err error = nil
	var sqlFilename string

	// Prolog.

	traceExitMessageNumber := 109
	if senzingSchema.isTrace {
		entryTime := time.Now()
		senzingSchema.traceEntry(100)
		defer func() { senzingSchema.traceExit(traceExitMessageNumber, err, time.Since(entryTime)) }()
	}

	// Determine which SQL file to process.

	parsedUrl, err := url.Parse(databaseUrl)
	if err != nil {
		if strings.HasPrefix(databaseUrl, "postgresql") {
			index := strings.LastIndex(databaseUrl, ":")
			newDatabaseUrl := databaseUrl[:index] + "/" + databaseUrl[index+1:]
			parsedUrl, err = url.Parse(newDatabaseUrl)
		}
		if err != nil {
			traceExitMessageNumber = 101
			return err
		}
	}

	switch parsedUrl.Scheme {
	case "sqlite3":
		sqlFilename = resourcePath + "/schema/g2core-schema-sqlite-create.sql"
	case "postgresql":
		sqlFilename = resourcePath + "/schema/g2core-schema-postgresql-create.sql"
	case "mysql":
		sqlFilename = resourcePath + "/schema/g2core-schema-mysql-create.sql"
	case "mssql":
		sqlFilename = resourcePath + "/schema/g2core-schema-mssql-create.sql"
	default:
		return fmt.Errorf("unknown database scheme: %s", parsedUrl.Scheme)
	}

	// Connect to the database.

	databaseConnector, err := connector.NewConnector(ctx, databaseUrl)
	if err != nil {
		traceExitMessageNumber = 102
		return err
	}

	// Create sqlExecutor to process file of SQL.

	sqlExecutor := &sqlexecutor.SqlExecutorImpl{
		DatabaseConnector: databaseConnector,
	}
	err = sqlExecutor.SetLogLevel(ctx, senzingSchema.logLevelName)
	if err != nil {
		traceExitMessageNumber = 103
		return err
	}

	// Add observers to sqlExecutor.

	if senzingSchema.observers != nil {
		for _, observer := range senzingSchema.observers.GetObservers(ctx) {
			err = sqlExecutor.RegisterObserver(ctx, observer)
			if err != nil {
				traceExitMessageNumber = 104
				return err
			}
		}
	}

	// Process file of SQL

	err = sqlExecutor.ProcessFileName(ctx, sqlFilename)
	if err != nil {
		traceExitMessageNumber = 105
		return err
	}
	senzingSchema.log(2001, sqlFilename, parsedUrl.Redacted())
	return err
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The InitializeSenzing method adds the Senzing database schema to the specified database.

Input
  - ctx: A context to control lifecycle.
*/
func (senzingSchema *SenzingSchemaImpl) InitializeSenzing(ctx context.Context) error {
	var err error = nil

	// Prolog.

	traceExitMessageNumber := 19
	if senzingSchema.isTrace {
		entryTime := time.Now()
		senzingSchema.traceEntry(10)
		defer func() { senzingSchema.traceExit(traceExitMessageNumber, err, time.Since(entryTime)) }()
	}

	// Log entry parameters.

	if senzingSchema.getLogger().IsDebug() {
		asJson, err := json.Marshal(senzingSchema)
		if err != nil {
			traceExitMessageNumber = 11
			return err
		}
		senzingSchema.log(1001, senzingSchema, string(asJson))
	}

	// Pull values out of SenzingEngineConfigurationJson.

	parser, err := engineconfigurationjsonparser.New(senzingSchema.SenzingEngineConfigurationJson)
	if err != nil {
		traceExitMessageNumber = 12
		return err
	}
	resourcePath, err := parser.GetResourcePath(ctx)
	if err != nil {
		traceExitMessageNumber = 13
		return err
	}
	databaseUrls, err := parser.GetDatabaseUrls(ctx)
	if err != nil {
		traceExitMessageNumber = 14
		return err
	}

	// Process each database.

	for _, databaseUrl := range databaseUrls {
		err = senzingSchema.processDatabase(ctx, resourcePath, databaseUrl)
		if err != nil {
			traceExitMessageNumber = 15
			return err
		}
	}

	// Notify observers.

	if senzingSchema.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, senzingSchema.observers, ProductId, 8001, err, details)
		}()
	}

	return err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (senzingSchema *SenzingSchemaImpl) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil

	// Prolog.

	traceExitMessageNumber := 29
	if senzingSchema.isTrace {
		entryTime := time.Now()
		senzingSchema.traceEntry(20, observer.GetObserverId(ctx))
		defer func() {
			senzingSchema.traceExit(traceExitMessageNumber, observer.GetObserverId(ctx), err, time.Since(entryTime))
		}()
	}

	// Log entry parameters.

	if senzingSchema.getLogger().IsDebug() {
		asJson, err := json.Marshal(senzingSchema)
		if err != nil {
			traceExitMessageNumber = 21
			return err
		}
		senzingSchema.log(1002, senzingSchema, string(asJson))
	}

	// Create empty list of observers.

	if senzingSchema.observers == nil {
		senzingSchema.observers = &subject.SubjectImpl{}
	}

	// Register observer with senzingSchema.

	err = senzingSchema.observers.RegisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber = 22
		return err
	}

	// Notify observers.

	go func() {
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, senzingSchema.observers, ProductId, 8002, err, details)
	}()

	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (senzingSchema *SenzingSchemaImpl) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil

	// Prolog.

	traceExitMessageNumber := 39
	if senzingSchema.isTrace {
		entryTime := time.Now()
		senzingSchema.traceEntry(30, logLevelName)
		defer func() { senzingSchema.traceExit(traceExitMessageNumber, logLevelName, err, time.Since(entryTime)) }()
	}

	// Log entry parameters.

	if senzingSchema.getLogger().IsDebug() {
		asJson, err := json.Marshal(senzingSchema)
		if err != nil {
			traceExitMessageNumber = 31
			return err
		}
		senzingSchema.log(1003, senzingSchema, string(asJson))
	}

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		traceExitMessageNumber = 32
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}

	// Set senzingConfig log level.

	senzingSchema.logLevelName = logLevelName
	senzingSchema.getLogger().SetLogLevel(logLevelName)
	senzingSchema.isTrace = (logLevelName == logging.LevelTraceName)

	// Notify observers.

	if senzingSchema.observers != nil { // Performance optimization.
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, senzingSchema.observers, ProductId, 8003, err, details)
		}()
	}

	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be removed.
*/
func (senzingSchema *SenzingSchemaImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil

	// Prolog.

	traceExitMessageNumber := 49
	if senzingSchema.isTrace {
		entryTime := time.Now()
		senzingSchema.traceEntry(40, observer.GetObserverId(ctx))
		defer func() {
			senzingSchema.traceExit(traceExitMessageNumber, observer.GetObserverId(ctx), err, time.Since(entryTime))
		}()
	}

	// Log entry parameters.

	if senzingSchema.getLogger().IsDebug() {
		asJson, err := json.Marshal(senzingSchema)
		if err != nil {
			traceExitMessageNumber = 41
			return err
		}
		senzingSchema.log(1004, senzingSchema, string(asJson))
	}

	// Remove observer from this service.

	if senzingSchema.observers != nil {

		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, senzingSchema.observers, ProductId, 8004, err, details)

		err = senzingSchema.observers.UnregisterObserver(ctx, observer)
		if err != nil {
			traceExitMessageNumber = 42
			return err
		}

		if !senzingSchema.observers.HasObservers(ctx) {
			senzingSchema.observers = nil
		}
	}

	return err
}
