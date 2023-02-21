package sqlfiler

import (
	"bufio"
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"os"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// InitializerImpl is the default implementation of the GrpcServer interface.
type SqlfilerImpl struct {
	isTrace           bool
	logger            messagelogger.MessageLoggerInterface
	LogLevel          logger.Level
	DatabaseConnector driver.Connector
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// func failOnError(msgId int, err error) {
// 	logger.Log(msgId, err)
// 	panic(err.Error())
// }

// Print error and leave program.
func errorExit(message string, err error) {
	fmt.Printf("Exit with error: %s   Error: %v\n", message, err)
	os.Exit(1)
}

// Get the Logger singleton.
func (sqlfiler *SqlfilerImpl) getLogger() messagelogger.MessageLoggerInterface {
	if sqlfiler.logger == nil {
		sqlfiler.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return sqlfiler.logger
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
func (sqlfiler *SqlfilerImpl) ProcessFileName(ctx context.Context, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	sqlfiler.ProcessScanner(ctx, bufio.NewScanner(file))
}

/*
The ProcessScanner does a database call for each line scanned.

Input
  - ctx: A context to control lifecycle.
  - scanner: SQL statements to be processed.
*/
func (sqlfiler *SqlfilerImpl) ProcessScanner(ctx context.Context, scanner *bufio.Scanner) {

	// Open a database connection.

	database := sql.OpenDB(sqlfiler.DatabaseConnector)
	defer database.Close()

	// Process each scanned line.

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		result, err := database.ExecContext(ctx, scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Result: %v", result)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
