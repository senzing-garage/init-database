package initializer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/senzing/go-common/engineconfigurationjsonparser"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-observing/notifier"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
	"github.com/senzing/init-database/senzingconfig"
	"github.com/senzing/init-database/senzingschema"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// InitializerImpl is the default implementation of the Initializer interface.
type InitializerImpl struct {
	DataSources                    []string
	logger                         logging.LoggingInterface
	observers                      subject.Subject
	senzingConfigSingleton         senzingconfig.SenzingConfig
	SenzingEngineConfigurationJson string
	SenzingLogLevel                string
	SenzingModuleName              string
	senzingSchemaSingleton         senzingschema.SenzingSchema
	SenzingVerboseLogging          int
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var debugOptions []interface{} = []interface{}{
	&logging.OptionCallerSkip{Value: 5},
}

var traceOptions []interface{} = []interface{}{
	&logging.OptionCallerSkip{Value: 5},
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (initializerImpl *InitializerImpl) getLogger() logging.LoggingInterface {
	var err error = nil
	if initializerImpl.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		initializerImpl.logger, err = logging.NewSenzingToolsLogger(ProductId, IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return initializerImpl.logger
}

// Log message.
func (initializerImpl *InitializerImpl) log(messageNumber int, details ...interface{}) {
	initializerImpl.getLogger().Log(messageNumber, details...)
}

// Debug.
func (initializerImpl *InitializerImpl) debug(messageNumber int, details ...interface{}) {
	details = append(details, debugOptions...)
	initializerImpl.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (initializerImpl *InitializerImpl) traceEntry(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	initializerImpl.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (initializerImpl *InitializerImpl) traceExit(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	initializerImpl.getLogger().Log(messageNumber, details...)
}

// --- Dependent services -----------------------------------------------------

func (initializerImpl *InitializerImpl) getSenzingConfig() senzingconfig.SenzingConfig {
	if initializerImpl.senzingConfigSingleton == nil {
		initializerImpl.senzingConfigSingleton = &senzingconfig.SenzingConfigImpl{
			DataSources:                    initializerImpl.DataSources,
			SenzingEngineConfigurationJson: initializerImpl.SenzingEngineConfigurationJson,
			SenzingModuleName:              initializerImpl.SenzingModuleName,
			SenzingVerboseLogging:          initializerImpl.SenzingVerboseLogging,
		}
	}
	return initializerImpl.senzingConfigSingleton
}

func (initializerImpl *InitializerImpl) getSenzingSchema() senzingschema.SenzingSchema {
	if initializerImpl.senzingSchemaSingleton == nil {
		initializerImpl.senzingSchemaSingleton = &senzingschema.SenzingSchemaImpl{
			SenzingEngineConfigurationJson: initializerImpl.SenzingEngineConfigurationJson,
		}
	}
	return initializerImpl.senzingSchemaSingleton
}

// --- Specific database processing -------------------------------------------

func (initializerImpl *InitializerImpl) initializeSpecificDatabaseSqlite(ctx context.Context, parsedUrl *url.URL) error {
	var err error = nil

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 109
	if initializerImpl.getLogger().IsDebug() {
		defer func() {
			if debugMessageNumber > 0 {
				initializerImpl.debug(debugMessageNumber, err)
			}
		}()
		if initializerImpl.getLogger().IsTrace() {
			entryTime := time.Now()
			initializerImpl.traceEntry(100)
			defer func() { initializerImpl.traceExit(traceExitMessageNumber, err, time.Since(entryTime)) }()
		}
	}

	// If file exists, no more to do.

	filename := parsedUrl.Path
	_, err = os.Stat(filename)
	if err == nil {
		return err // Nothing more to do.
	}

	// File doesn't exist, create it.

	path := filepath.Dir(filename)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		debugMessageNumber = 1101
		traceExitMessageNumber = 101
		return err
	}
	_, err = os.Create(filename)
	if err != nil {
		debugMessageNumber = 1102
		traceExitMessageNumber = 102
		return err
	}
	initializerImpl.log(2001, filename)

	// Notify observers.

	if initializerImpl.observers != nil {
		go func() {
			details := map[string]string{
				"sqliteFile": filename,
			}
			notifier.Notify(ctx, initializerImpl.observers, ProductId, 8005, err, details)
		}()
	}
	return err
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

	// Initialize logging.

	logLevel := initializerImpl.SenzingLogLevel
	if logLevel == "" {
		logLevel = "INFO"
	}
	err = initializerImpl.SetLogLevel(ctx, logLevel)
	if err != nil {
		return err
	}

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 19
	if initializerImpl.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializerImpl.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializerImpl.getLogger().IsTrace() {
			entryTime := time.Now()
			initializerImpl.traceEntry(10)
			defer func() { initializerImpl.traceExit(traceExitMessageNumber, err, time.Since(entryTime)) }()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJson, err := json.Marshal(initializerImpl)
		if err != nil {
			debugMessageNumber = 1011
			traceExitMessageNumber = 11
			return err
		}
		initializerImpl.log(1000, initializerImpl, string(asJson))
	}

	// Perform initialization for specific databases.

	err = initializerImpl.InitializeSpecificDatabase(ctx)
	if err != nil {
		debugMessageNumber = 1012
		traceExitMessageNumber = 12
		return err
	}

	// Create schema in database.

	senzingSchema := initializerImpl.getSenzingSchema()
	err = senzingSchema.SetLogLevel(ctx, logLevel)
	if err != nil {
		debugMessageNumber = 1013
		traceExitMessageNumber = 13
		return err
	}
	err = senzingSchema.InitializeSenzing(ctx)
	if err != nil {
		debugMessageNumber = 1014
		traceExitMessageNumber = 14
		return err
	}

	// Create initial Senzing configuration.

	senzingConfig := initializerImpl.getSenzingConfig()
	senzingConfig.SetLogLevel(ctx, logLevel)
	if err != nil {
		debugMessageNumber = 1015
		traceExitMessageNumber = 15
		return err
	}
	err = senzingConfig.InitializeSenzing(ctx)
	if err != nil {
		debugMessageNumber = 1016
		traceExitMessageNumber = 16
		return err
	}

	// Notify observers.

	if initializerImpl.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, initializerImpl.observers, ProductId, 8001, err, details)
		}()
	}
	return err
}

/*
The InitializeSpecificDatabase method routes specific databse processing
based on the database URL's protocol field.

Input
  - ctx: A context to control lifecycle.
*/
func (initializerImpl *InitializerImpl) InitializeSpecificDatabase(ctx context.Context) error {
	var err error = nil
	var databaseUrls []string

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 29
	if initializerImpl.getLogger().IsDebug() {
		defer func() {
			if debugMessageNumber > 0 {
				initializerImpl.debug(debugMessageNumber, err)
			}
		}()
		if initializerImpl.getLogger().IsTrace() {
			entryTime := time.Now()
			initializerImpl.traceEntry(20)
			defer func() { initializerImpl.traceExit(traceExitMessageNumber, err, time.Since(entryTime)) }()
		}
		asJson, err := json.Marshal(initializerImpl)
		if err != nil {
			debugMessageNumber = 1021
			traceExitMessageNumber = 21
			return err
		}
		initializerImpl.log(1001, initializerImpl, string(asJson))
	}

	// Pull values out of SenzingEngineConfigurationJson.

	parser, err := engineconfigurationjsonparser.New(initializerImpl.SenzingEngineConfigurationJson)
	if err != nil {
		debugMessageNumber = 1022
		traceExitMessageNumber = 22
		return err
	}
	databaseUrls, err = parser.GetDatabaseUrls(ctx)
	if err != nil {
		debugMessageNumber = 1023
		traceExitMessageNumber = 23
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
				debugMessageNumber = 1024
				traceExitMessageNumber = 24
				return err
			}
		}

		// Special handling for each database type.

		switch parsedUrl.Scheme {
		case "sqlite3":
			err = initializerImpl.initializeSpecificDatabaseSqlite(ctx, parsedUrl)
			if err != nil {
				debugMessageNumber = 1025
				traceExitMessageNumber = 25
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
	var err error = nil

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 39
	if initializerImpl.getLogger().IsDebug() {
		defer func() {
			if debugMessageNumber > 0 {
				initializerImpl.debug(debugMessageNumber, err)
			}
		}()
		if initializerImpl.getLogger().IsTrace() {
			entryTime := time.Now()
			initializerImpl.traceEntry(30)
			defer func() {
				initializerImpl.traceExit(traceExitMessageNumber, observer.GetObserverId(ctx), err, time.Since(entryTime))
			}()
		}
		asJson, err := json.Marshal(initializerImpl)
		if err != nil {
			debugMessageNumber = 1031
			traceExitMessageNumber = 31
			return err
		}
		initializerImpl.log(1002, initializerImpl, string(asJson))
	}

	// Create empty list of observers.

	if initializerImpl.observers == nil {
		initializerImpl.observers = &subject.SubjectImpl{}
	}

	// Register observer with initializerImpl and dependent services.

	err = initializerImpl.observers.RegisterObserver(ctx, observer)
	if err != nil {
		debugMessageNumber = 1032
		traceExitMessageNumber = 32
		return err
	}
	err = initializerImpl.getSenzingConfig().RegisterObserver(ctx, observer)
	if err != nil {
		debugMessageNumber = 1033
		traceExitMessageNumber = 33
		return err
	}
	err = initializerImpl.getSenzingSchema().RegisterObserver(ctx, observer)
	if err != nil {
		debugMessageNumber = 1034
		traceExitMessageNumber = 34
		return err
	}

	// Notify observers.

	go func() {
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, initializerImpl.observers, ProductId, 8002, err, details)
	}()
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (initializerImpl *InitializerImpl) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 49
	if initializerImpl.getLogger().IsDebug() {
		defer func() {
			if debugMessageNumber > 0 {
				initializerImpl.debug(debugMessageNumber, err)
			}
		}()
		if initializerImpl.getLogger().IsTrace() {
			entryTime := time.Now()
			initializerImpl.traceEntry(40)
			defer func() { initializerImpl.traceExit(traceExitMessageNumber, logLevelName, err, time.Since(entryTime)) }()
		}
		asJson, err := json.Marshal(initializerImpl)
		if err != nil {
			debugMessageNumber = 1041
			traceExitMessageNumber = 41
			return err
		}
		initializerImpl.log(1003, initializerImpl, string(asJson))
	}

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		debugMessageNumber = 1042
		traceExitMessageNumber = 42
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}

	// Set initializerImpl log level.

	err = initializerImpl.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		debugMessageNumber = 1043
		traceExitMessageNumber = 43
		return err
	}

	// Set log level for dependent services.

	if initializerImpl.senzingConfigSingleton != nil {
		err = initializerImpl.senzingConfigSingleton.SetLogLevel(ctx, logLevelName)
		if err != nil {
			debugMessageNumber = 1044
			traceExitMessageNumber = 44
			return err
		}
	}
	err = initializerImpl.getSenzingSchema().SetLogLevel(ctx, logLevelName)
	if err != nil {
		debugMessageNumber = 1045
		traceExitMessageNumber = 45
		return err
	}

	// Notify observers.

	if initializerImpl.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, initializerImpl.observers, ProductId, 8003, err, details)
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
func (initializerImpl *InitializerImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 59
	if initializerImpl.getLogger().IsDebug() {
		defer func() {
			if debugMessageNumber > 0 {
				initializerImpl.debug(debugMessageNumber, err)
			}
		}()
		if initializerImpl.getLogger().IsTrace() {
			entryTime := time.Now()
			initializerImpl.traceEntry(50)
			defer func() {
				initializerImpl.traceExit(traceExitMessageNumber, observer.GetObserverId(ctx), err, time.Since(entryTime))
			}()
		}
		asJson, err := json.Marshal(initializerImpl)
		if err != nil {
			debugMessageNumber = 1051
			traceExitMessageNumber = 51
			return err
		}
		initializerImpl.log(1004, initializerImpl, string(asJson))
	}

	// Unregister observer in dependencies.

	// err = initializerImpl.getSenzingConfig().UnregisterObserver(ctx, observer)
	// if err != nil {
	//	debugMessageNumber = 1052
	// 	traceExitMessageNumber = 52
	// 	return err
	// }
	// err = initializerImpl.getSenzingSchema().UnregisterObserver(ctx, observer)
	// if err != nil {
	//  debugMessageNumber = 1053
	// 	traceExitMessageNumber = 53
	// 	return err
	// }

	// Remove observer from this service.

	if initializerImpl.observers != nil {

		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, initializerImpl.observers, ProductId, 8004, err, details)

		err = initializerImpl.observers.UnregisterObserver(ctx, observer)
		if err != nil {
			debugMessageNumber = 1054
			traceExitMessageNumber = 54
			return err
		}

		if !initializerImpl.observers.HasObservers(ctx) {
			initializerImpl.observers = nil
		}
	}
	return err
}
