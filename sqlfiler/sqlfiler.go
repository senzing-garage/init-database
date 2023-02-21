package sqlfiler

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"database/sql/driver"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// InitializerImpl is the default implementation of the GrpcServer interface.
type SqlfilerImpl struct {
	isTrace   bool
	logger    messagelogger.MessageLoggerInterface
	LogLevel  logger.Level
	connector driver.Connector
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

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

func (sqlfiler *SqlfilerImpl) ProcessFileName(ctx context.Context, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sqlfiler.ProcessScanner(ctx, bufio.NewScanner(file))
}

func (sqlfiler *SqlfilerImpl) ProcessScanner(ctx context.Context, scanner *bufio.Scanner) {

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
