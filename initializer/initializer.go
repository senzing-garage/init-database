package initializer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/senzing/go-common/engineconfigurationjsonparser"
	"github.com/senzing/go-databasing/connector"
	"github.com/senzing/go-databasing/sqlexecutor"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
	"github.com/senzing/go-sdk-abstract-factory/factory"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// InitializerImpl is the default implementation of the GrpcServer interface.
type InitializerImpl struct {
	GrpcAddress                    string
	GrpcOptions                    []grpc.DialOption
	isTrace                        bool
	messageLogger                  messagelogger.MessageLoggerInterface
	LogLevel                       logger.Level
	observers                      subject.Subject
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int
	g2configSingleton              g2api.G2config
	g2configSyncOnce               sync.Once
	g2configmgrSingleton           g2api.G2configmgr
	g2configmgrSyncOnce            sync.Once
	g2factorySingleton             factory.SdkAbstractFactory
	g2factorySyncOnce              sync.Once
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var defaultModuleName string = "initdatabase"

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Print error and leave program.
func errorExit(message string, err error) {
	fmt.Printf("Exit with error: %s   Error: %v\n", message, err)
	os.Exit(1)
}

// Get the Logger singleton.
func (initializer *InitializerImpl) getLogger() messagelogger.MessageLoggerInterface {
	if initializer.messageLogger == nil {
		initializer.messageLogger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return initializer.messageLogger
}

// Notify registered observers.
func (initializer *InitializerImpl) notify(ctx context.Context, messageId int, err error, details map[string]string) {
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
		initializer.observers.NotifyObservers(ctx, string(message))
	}
}

// Trace method entry.
func (initializer *InitializerImpl) traceEntry(errorNumber int, details ...interface{}) {
	initializer.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (initializer *InitializerImpl) traceExit(errorNumber int, details ...interface{}) {
	initializer.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Initialize method adds the Senzing database schema and Senzing default configuration to databases.

Input
  - ctx: A context to control lifecycle.
*/
func (initializer *InitializerImpl) InitializeSenzingSchema(ctx context.Context) error {
	if initializer.isTrace {
		initializer.traceEntry(3)
	}
	entryTime := time.Now()

	// Log entry parameters.

	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, initializer.LogLevel)
	logger.Log(2000, initializer)

	parser, err := engineconfigurationjsonparser.New(initializer.SenzingEngineConfigurationJson)
	if err != nil {
		return err
	}

	// Pull values out of SENZING_ENGINE_CONFIGURATION_JSON.

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

		// Process file of SQL.

		sqlExecutor := &sqlexecutor.SqlExecutorImpl{
			DatabaseConnector: databaseConnector,
		}
		// FIXME: sqlExecutor.SetLogLevel(ctx, (logger.LogLevel)initializer.messageLogger.GetLogLevel())
		// FIXME: When available, use initializer.observers.RegisterObservers()

		err = sqlExecutor.ProcessFileName(ctx, sqlFilename)
		if err != nil {
			return err
		}
	}

	if initializer.observers != nil {
		go func() {
			details := map[string]string{}
			initializer.notify(ctx, 8003, err, details)
		}()
	}
	if initializer.isTrace {
		defer initializer.traceExit(4, err, time.Since(entryTime))
	}
	return err
}

/*
The Initialize method adds the Senzing database schema and Senzing default configuration to databases.
Essentially it calls InitializeSenzingConfiguration() and InitializeSenzingSchema().

Input
  - ctx: A context to control lifecycle.
*/
func (initializer *InitializerImpl) Initialize(ctx context.Context) error {
	if initializer.isTrace {
		initializer.traceEntry(5)
	}
	entryTime := time.Now()

	// Initialize schema and configuration.

	err := initializer.InitializeSenzingSchema(ctx)
	if err != nil {
		errorExit("Could not initialize Senzing database schema.", err)
	}
	err = initializer.InitializeSenzingConfiguration(ctx)
	if err != nil {
		errorExit("Could not create Senzing configuration.", err)
	}

	// Epilog.

	if initializer.observers != nil {
		go func() {
			details := map[string]string{}
			initializer.notify(ctx, 8004, err, details)
		}()
	}
	if initializer.isTrace {
		defer initializer.traceExit(6, err, time.Since(entryTime))
	}
	return err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (initializer *InitializerImpl) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	if initializer.isTrace {
		initializer.traceEntry(7, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	if initializer.observers == nil {
		initializer.observers = &subject.SubjectImpl{}
	}
	err := initializer.observers.RegisterObserver(ctx, observer)
	if initializer.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverId(ctx),
			}
			initializer.notify(ctx, 8005, err, details)
		}()
	}
	if initializer.isTrace {
		defer initializer.traceExit(8, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (initializer *InitializerImpl) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if initializer.isTrace {
		initializer.traceEntry(9, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	initializer.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	initializer.isTrace = (initializer.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if initializer.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logger.LevelToTextMap[logLevel],
			}
			initializer.notify(ctx, 8006, err, details)
		}()
	}
	if initializer.isTrace {
		defer initializer.traceExit(10, logLevel, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be removed.
*/
func (initializer *InitializerImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	if initializer.isTrace {
		initializer.traceEntry(11, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	var err error = nil
	if initializer.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		initializer.notify(ctx, 8007, err, details)
	}
	err = initializer.observers.UnregisterObserver(ctx, observer)
	if !initializer.observers.HasObservers(ctx) {
		initializer.observers = nil
	}
	if initializer.isTrace {
		defer initializer.traceExit(12, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}
