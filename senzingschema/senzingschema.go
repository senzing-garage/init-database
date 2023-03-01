package senzingschema

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/senzing/go-common/engineconfigurationjsonparser"
	"github.com/senzing/go-databasing/connector"
	"github.com/senzing/go-databasing/sqlexecutor"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
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
	messageLogger                  messagelogger.MessageLoggerInterface
	logLevel                       logger.Level
	observers                      subject.Subject
	SenzingEngineConfigurationJson string
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (senzingSchema *SenzingSchemaImpl) getLogger() messagelogger.MessageLoggerInterface {
	if senzingSchema.messageLogger == nil {
		senzingSchema.messageLogger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return senzingSchema.messageLogger
}

// Trace method entry.
func (senzingSchema *SenzingSchemaImpl) traceEntry(errorNumber int, details ...interface{}) {
	senzingSchema.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (senzingSchema *SenzingSchemaImpl) traceExit(errorNumber int, details ...interface{}) {
	senzingSchema.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Initialize method adds the Senzing database schema to the specified database.

Input
  - ctx: A context to control lifecycle.
*/
func (senzingSchema *SenzingSchemaImpl) Initialize(ctx context.Context) error {
	if senzingSchema.isTrace {
		senzingSchema.traceEntry(1)
	}
	entryTime := time.Now()

	// Log entry parameters.

	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, senzingSchema.logLevel)
	logger.Log(1000, senzingSchema)

	// Pull values out of SENZING_ENGINE_CONFIGURATION_JSON.

	parser, err := engineconfigurationjsonparser.New(senzingSchema.SenzingEngineConfigurationJson)
	if err != nil {
		return err
	}
	resourcePath, err := parser.GetResourcePath(ctx)
	if err != nil {
		return err
	}
	databaseUrls, err := parser.GetDatabaseUrls(ctx)
	if err != nil {
		return err
	}

	// Process each database.

	for _, databaseUrl := range databaseUrls {
		var sqlFilename string

		// Connect to the database.

		databaseConnector, err := connector.NewConnector(ctx, databaseUrl)
		if err != nil {
			return err
		}

		// Determine which SQL file to process.

		parsedUrl, err := url.Parse(databaseUrl)
		if err != nil {
			return err
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

		// Create sqlExecutor to process file of SQL.

		sqlExecutor := &sqlexecutor.SqlExecutorImpl{
			DatabaseConnector: databaseConnector,
		}
		sqlExecutor.SetLogLevel(ctx, senzingSchema.logLevel)

		// Add observers to sqlExecutor.

		if senzingSchema.observers != nil {
			for _, observer := range senzingSchema.observers.GetObservers(ctx) {
				err = sqlExecutor.RegisterObserver(ctx, observer)
				if err != nil {
					return err
				}
			}
		}

		// Process file of SQL

		err = sqlExecutor.ProcessFileName(ctx, sqlFilename)
		if err != nil {
			return err
		}
	}

	// Epilog.

	if senzingSchema.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, senzingSchema.observers, ProductId, 8001, err, details)
		}()
	}
	if senzingSchema.isTrace {
		defer senzingSchema.traceExit(2, err, time.Since(entryTime))
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
	if senzingSchema.isTrace {
		senzingSchema.traceEntry(3, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	if senzingSchema.observers == nil {
		senzingSchema.observers = &subject.SubjectImpl{}
	}
	err := senzingSchema.observers.RegisterObserver(ctx, observer)
	details := map[string]string{
		"observerID": observer.GetObserverId(ctx),
	}
	notifier.Notify(ctx, senzingSchema.observers, ProductId, 8002, err, details)
	if senzingSchema.isTrace {
		defer senzingSchema.traceExit(4, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (senzingSchema *SenzingSchemaImpl) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if senzingSchema.isTrace {
		senzingSchema.traceEntry(5, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	senzingSchema.logLevel = logLevel
	senzingSchema.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	senzingSchema.isTrace = (senzingSchema.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if senzingSchema.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logger.LevelToTextMap[logLevel],
			}
			notifier.Notify(ctx, senzingSchema.observers, ProductId, 8003, err, details)
		}()
	}
	if senzingSchema.isTrace {
		defer senzingSchema.traceExit(6, logLevel, err, time.Since(entryTime))
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
	if senzingSchema.isTrace {
		senzingSchema.traceEntry(7, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	var err error = nil
	if senzingSchema.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, senzingSchema.observers, ProductId, 8004, err, details)
	}
	err = senzingSchema.observers.UnregisterObserver(ctx, observer)
	if !senzingSchema.observers.HasObservers(ctx) {
		senzingSchema.observers = nil
	}
	if senzingSchema.isTrace {
		defer senzingSchema.traceExit(8, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}
