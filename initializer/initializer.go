package initializer

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/senzing/go-common/engineconfigurationjsonparser"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/notifier"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
	"github.com/senzing/initdatabase/senzingconfig"
	"github.com/senzing/initdatabase/senzingschema"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// InitializerImpl is the default implementation of the Initializer interface.
type InitializerImpl struct {
	DataSources                    []string
	isTrace                        bool
	logLevel                       logger.Level
	messageLogger                  messagelogger.MessageLoggerInterface
	observers                      subject.Subject
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (initializerImpl *InitializerImpl) getLogger() messagelogger.MessageLoggerInterface {
	if initializerImpl.messageLogger == nil {
		initializerImpl.messageLogger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, initializerImpl.logLevel)
	}
	return initializerImpl.messageLogger
}

func (initializerImpl *InitializerImpl) initializeSpecificDatabaseSqlite(ctx context.Context, parsedUrl *url.URL) error {
	// If file doesn't exist, create it.

	filename := parsedUrl.Path
	_, err := os.Stat(filename)
	if err != nil {
		path := filepath.Dir(filename)
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		_, err = os.Create(filename)
		if err != nil {
			return err
		}
		initializerImpl.getLogger().Log(2001, filename)
		if initializerImpl.observers != nil {
			go func() {
				details := map[string]string{
					"sqliteFile": filename,
				}
				notifier.Notify(ctx, initializerImpl.observers, ProductId, 8005, err, details)
			}()
		}
	}
	return err
}

// Trace method entry.
func (initializerImpl *InitializerImpl) traceEntry(errorNumber int, details ...interface{}) {
	initializerImpl.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (initializerImpl *InitializerImpl) traceExit(errorNumber int, details ...interface{}) {
	initializerImpl.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Initialize method adds the Senzing database schema and Senzing default configuration to databases.
Essentially it calls senzingSchema.Initialize() and senzingConfig.Initialize(ctx).

Input
  - ctx: A context to control lifecycle.
*/
func (initializerImpl *InitializerImpl) Initialize(ctx context.Context) error {
	var err error = nil
	if initializerImpl.isTrace {
		initializerImpl.traceEntry(1)
	}
	entryTime := time.Now()

	// Log entry parameters.

	initializerImpl.getLogger().Log(1000, initializerImpl)

	// Create senzingSchema for initializing Senzing schema.

	senzingSchema := &senzingschema.SenzingSchemaImpl{
		SenzingEngineConfigurationJson: initializerImpl.SenzingEngineConfigurationJson,
	}
	senzingSchema.SetLogLevel(ctx, initializerImpl.logLevel)

	// Create senzingConfig for initializing Senzing configuration.

	senzingConfig := &senzingconfig.SenzingConfigImpl{
		DataSources:                    initializerImpl.DataSources,
		SenzingEngineConfigurationJson: initializerImpl.SenzingEngineConfigurationJson,
		SenzingModuleName:              initializerImpl.SenzingModuleName,
		SenzingVerboseLogging:          initializerImpl.SenzingVerboseLogging,
	}
	senzingConfig.SetLogLevel(ctx, initializerImpl.logLevel)

	// Add observers to structs.

	if initializerImpl.observers != nil {
		for _, observer := range initializerImpl.observers.GetObservers(ctx) {
			err = senzingConfig.RegisterObserver(ctx, observer)
			if err != nil {
				return err
			}
			err = senzingSchema.RegisterObserver(ctx, observer)
			if err != nil {
				return err
			}
		}
	}

	// Perform initialization.

	err = initializerImpl.InitializeSpecificDatabase(ctx)
	if err != nil {
		return err
	}
	err = senzingSchema.Initialize(ctx)
	if err != nil {
		return err
	}
	err = senzingConfig.Initialize(ctx)
	if err != nil {
		return err
	}

	// Epilog.

	if initializerImpl.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, initializerImpl.observers, ProductId, 8001, err, details)
		}()
	}
	if initializerImpl.isTrace {
		defer initializerImpl.traceExit(2, err, time.Since(entryTime))
	}
	return err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
*/
func (initializerImpl *InitializerImpl) InitializeSpecificDatabase(ctx context.Context) error {

	// Pull values out of SenzingEngineConfigurationJson.

	parser, err := engineconfigurationjsonparser.New(initializerImpl.SenzingEngineConfigurationJson)
	if err != nil {
		return err
	}
	databaseUrls, err := parser.GetDatabaseUrls(ctx)
	if err != nil {
		return err
	}

	// Process each database.

	for _, databaseUrl := range databaseUrls {

		// Parse URL.

		parsedUrl, err := url.Parse(databaseUrl)
		if err != nil {
			if strings.HasPrefix(databaseUrl, "postgresql") {
				index := strings.LastIndex(databaseUrl, ":")
				newDatabaseUrl := databaseUrl[:index] + "/" + databaseUrl[index+1:]
				parsedUrl, err = url.Parse(newDatabaseUrl)
			}
			if err != nil {
				return err
			}
		}

		// Special handling for each database type.

		switch parsedUrl.Scheme {
		case "sqlite3":
			err = initializerImpl.initializeSpecificDatabaseSqlite(ctx, parsedUrl)
			if err != nil {
				return err
			}
		}
	}
	return err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (initializerImpl *InitializerImpl) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	if initializerImpl.isTrace {
		initializerImpl.traceEntry(3, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	if initializerImpl.observers == nil {
		initializerImpl.observers = &subject.SubjectImpl{}
	}
	err := initializerImpl.observers.RegisterObserver(ctx, observer)
	details := map[string]string{
		"observerID": observer.GetObserverId(ctx),
	}
	notifier.Notify(ctx, initializerImpl.observers, ProductId, 8002, err, details)
	if initializerImpl.isTrace {
		defer initializerImpl.traceExit(4, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (initializerImpl *InitializerImpl) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if initializerImpl.isTrace {
		initializerImpl.traceEntry(5, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	initializerImpl.logLevel = logLevel
	initializerImpl.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	initializerImpl.isTrace = (initializerImpl.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if initializerImpl.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logger.LevelToTextMap[logLevel],
			}
			notifier.Notify(ctx, initializerImpl.observers, ProductId, 8003, err, details)
		}()
	}
	if initializerImpl.isTrace {
		defer initializerImpl.traceExit(6, logLevel, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be removed.
*/
func (initializerImpl *InitializerImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	if initializerImpl.isTrace {
		initializerImpl.traceEntry(7, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	var err error = nil
	if initializerImpl.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, initializerImpl.observers, ProductId, 8004, err, details)
	}
	err = initializerImpl.observers.UnregisterObserver(ctx, observer)
	if !initializerImpl.observers.HasObservers(ctx) {
		initializerImpl.observers = nil
	}
	if initializerImpl.isTrace {
		defer initializerImpl.traceExit(8, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}
