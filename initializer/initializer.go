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

	"github.com/senzing-garage/go-helpers/settingsparser"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/observerpb"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing-garage/init-database/senzingconfig"
	"github.com/senzing-garage/init-database/senzingschema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// InitializerImpl is the default implementation of the Initializer interface.
type InitializerImpl struct {
	DataSources                    []string
	logger                         logging.Logging
	ObserverOrigin                 string
	observers                      subject.Subject
	ObserverUrl                    string
	senzingConfigSingleton         senzingconfig.SenzingConfig
	SenzingEngineConfigurationFile string
	SenzingEngineConfigurationJson string
	SenzingLogLevel                string
	SenzingModuleName              string
	senzingSchemaSingleton         senzingschema.SenzingSchema
	SenzingVerboseLogging          int64
	SqlFile                        string
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
func (initializerImpl *InitializerImpl) getLogger() logging.Logging {
	var err error = nil
	if initializerImpl.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		initializerImpl.logger, err = logging.NewSenzingLogger(ComponentId, IdMessages, options...)
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

// --- Observing --------------------------------------------------------------

func (initializerImpl *InitializerImpl) createGrpcObserver(ctx context.Context, parsedUrl url.URL) (observer.Observer, error) {
	var err error
	var result observer.Observer

	port := DefaultGrpcObserverPort
	if len(parsedUrl.Port()) > 0 {
		port = parsedUrl.Port()
	}
	target := fmt.Sprintf("%s:%s", parsedUrl.Hostname(), port)

	// TODO: Allow specification of options from ObserverUrl/parsedUrl
	grpcOptions := grpc.WithTransportCredentials(insecure.NewCredentials())

	grpcConnection, err := grpc.Dial(target, grpcOptions)
	if err != nil {
		return result, err
	}
	result = &observer.GrpcObserver{
		GrpcClient: observerpb.NewObserverClient(grpcConnection),
		ID:         "init-database",
	}
	return result, err
}

func (initializerImpl *InitializerImpl) registerObserverLocal(ctx context.Context, observer observer.Observer) error {
	if initializerImpl.observers == nil {
		initializerImpl.observers = &subject.SimpleSubject{}
	}
	return initializerImpl.observers.RegisterObserver(ctx, observer)
}

func (initializerImpl *InitializerImpl) registerObserverSenzingConfig(ctx context.Context, observer observer.Observer) error {
	initializerImpl.getSenzingConfig().SetObserverOrigin(ctx, initializerImpl.ObserverOrigin)
	return initializerImpl.getSenzingConfig().RegisterObserver(ctx, observer)
}

func (initializerImpl *InitializerImpl) registerObserverSenzingSchema(ctx context.Context, observer observer.Observer) error {
	initializerImpl.getSenzingSchema().SetObserverOrigin(ctx, initializerImpl.ObserverOrigin)
	return initializerImpl.getSenzingSchema().RegisterObserver(ctx, observer)
}

// --- Dependent services -----------------------------------------------------

func (initializerImpl *InitializerImpl) getSenzingConfig() senzingconfig.SenzingConfig {
	if initializerImpl.senzingConfigSingleton == nil {
		initializerImpl.senzingConfigSingleton = &senzingconfig.SenzingConfigImpl{
			DataSources:                    initializerImpl.DataSources,
			SenzingEngineConfigurationFile: initializerImpl.SenzingEngineConfigurationFile,
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
			SqlFile:                        initializerImpl.SqlFile,
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

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializerImpl.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializerImpl.getLogger().IsTrace() {
			entryTime := time.Now()
			initializerImpl.traceEntry(100, parsedUrl)
			defer func() { initializerImpl.traceExit(traceExitMessageNumber, parsedUrl, err, time.Since(entryTime)) }()
		}
	}

	// If file exists, no more to do.

	filename := parsedUrl.Path
	filename = filepath.Clean(filename) // See https://securego.io/docs/rules/g304.html
	_, err = os.Stat(filename)
	if err == nil {
		traceExitMessageNumber, debugMessageNumber = 101, 0 // debugMessageNumber=0 because it's not an error.
		return err                                          // Nothing more to do.
	}

	// File doesn't exist, create it.

	path := filepath.Dir(filename)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 102, 1102
		return err
	}
	_, err = os.Create(filename)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 103, 1103
		return err
	}
	initializerImpl.log(2001, filename)

	// Notify observers.

	if initializerImpl.observers != nil {
		go func() {
			details := map[string]string{
				"sqliteFile": filename,
			}
			notifier.Notify(ctx, initializerImpl.observers, initializerImpl.ObserverOrigin, ComponentId, 8010, err, details)
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
	debugMessageNumber := 0
	traceExitMessageNumber := 19

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
			traceExitMessageNumber, debugMessageNumber = 11, 1011
			return err
		}
		initializerImpl.log(1000, initializerImpl, string(asJson))
	}

	// Initialize observing.

	var anObserver observer.Observer
	if len(initializerImpl.ObserverUrl) > 0 {
		parsedUrl, err := url.Parse(initializerImpl.ObserverUrl)
		if err != nil {
			return err
		}
		switch parsedUrl.Scheme {
		case "grpc":
			anObserver, err = initializerImpl.createGrpcObserver(ctx, *parsedUrl)
			if err != nil {
				traceExitMessageNumber, debugMessageNumber = 18, 1018
				return err
			}
		}
		err = initializerImpl.registerObserverLocal(ctx, anObserver)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 17, 1017
			return err
		}

		// Notify observers.

		go func() {
			details := map[string]string{
				"observerID": anObserver.GetObserverID(ctx),
			}
			notifier.Notify(ctx, initializerImpl.observers, initializerImpl.ObserverOrigin, ComponentId, 8001, err, details)
		}()
	}

	// Verify database file exists.

	if len(initializerImpl.SqlFile) > 0 {
		_, err = os.Stat(initializerImpl.SqlFile)
		if err != nil {
			initializerImpl.log(3001, initializerImpl.SqlFile)
			traceExitMessageNumber, debugMessageNumber = 21, 1075
			return err
		}
	}

	// Perform initialization for specific databases.

	err = initializerImpl.InitializeSpecificDatabase(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 12, 1012
		return err
	}

	// Create schema in database.

	senzingSchema := initializerImpl.getSenzingSchema()
	err = senzingSchema.SetLogLevel(ctx, logLevel)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 13, 1013
		return err
	}
	err = initializerImpl.registerObserverSenzingSchema(ctx, anObserver)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 19, 1019
		return err
	}
	err = senzingSchema.InitializeSenzing(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 14, 1014
		return err
	}

	// Create initial Senzing configuration.

	senzingConfig := initializerImpl.getSenzingConfig()
	err = senzingConfig.SetLogLevel(ctx, logLevel)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 15, 1015
		return err
	}
	err = initializerImpl.registerObserverSenzingConfig(ctx, anObserver)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 20, 1000
		return err
	}
	err = senzingConfig.InitializeSenzing(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 16, 1016
		return err
	}

	// Notify observers.

	if initializerImpl.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, initializerImpl.observers, initializerImpl.ObserverOrigin, ComponentId, 8002, err, details)
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
	traceExitMessageNumber := 49
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
			initializerImpl.traceEntry(40)
			defer func() { initializerImpl.traceExit(traceExitMessageNumber, err, time.Since(entryTime)) }()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJson, err := json.Marshal(initializerImpl)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 41, 1041
			return err
		}
		initializerImpl.log(1001, initializerImpl, string(asJson))
	}

	// Pull values out of SenzingEngineConfigurationJson.

	parser, err := settingsparser.New(initializerImpl.SenzingEngineConfigurationJson)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 42, 1042
		return err
	}
	databaseUrls, err = parser.GetDatabaseURLs(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 43, 1043
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
				traceExitMessageNumber, debugMessageNumber = 44, 1044
				return err
			}
		}

		// Special handling for each database type.

		switch parsedUrl.Scheme {
		case "sqlite3":
			err = initializerImpl.initializeSpecificDatabaseSqlite(ctx, parsedUrl)
			if err != nil {
				traceExitMessageNumber, debugMessageNumber = 45, 1045
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

	if observer == nil {
		return err
	}

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 59
	if initializerImpl.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializerImpl.debug(debugMessageNumber, observer.GetObserverID(ctx), err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializerImpl.getLogger().IsTrace() {
			entryTime := time.Now()
			initializerImpl.traceEntry(50, observer.GetObserverID(ctx))
			defer func() {
				initializerImpl.traceExit(traceExitMessageNumber, observer.GetObserverID(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJson, err := json.Marshal(initializerImpl)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 51, 1051
			return err
		}
		initializerImpl.log(1002, initializerImpl, string(asJson))
	}

	// Create empty list of observers.

	if initializerImpl.observers == nil {
		initializerImpl.observers = &subject.SimpleSubject{}
	}

	// Register observer with initializerImpl and dependent services.

	err = initializerImpl.observers.RegisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 52, 1052
		return err
	}
	err = initializerImpl.getSenzingConfig().RegisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 53, 1053
		return err
	}
	err = initializerImpl.getSenzingSchema().RegisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 54, 1054
		return err
	}

	// Notify observers.

	go func() {
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, initializerImpl.observers, initializerImpl.ObserverOrigin, ComponentId, 8003, err, details)
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
	traceExitMessageNumber := 69
	if initializerImpl.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializerImpl.debug(debugMessageNumber, logLevelName, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializerImpl.getLogger().IsTrace() {
			entryTime := time.Now()
			initializerImpl.traceEntry(60, logLevelName)
			defer func() { initializerImpl.traceExit(traceExitMessageNumber, logLevelName, err, time.Since(entryTime)) }()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJson, err := json.Marshal(initializerImpl)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 61, 1061
			return err
		}
		initializerImpl.log(1003, initializerImpl, string(asJson))
	}

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		traceExitMessageNumber, debugMessageNumber = 62, 1062
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}

	// Set initializerImpl log level.

	err = initializerImpl.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 63, 1063
		return err
	}

	// Set log level for dependent services.

	if initializerImpl.senzingConfigSingleton != nil {
		err = initializerImpl.senzingConfigSingleton.SetLogLevel(ctx, logLevelName)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 64, 1064
			return err
		}
	}
	err = initializerImpl.getSenzingSchema().SetLogLevel(ctx, logLevelName)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 65, 1065
		return err
	}

	// Notify observers.

	if initializerImpl.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, initializerImpl.observers, initializerImpl.ObserverOrigin, ComponentId, 8004, err, details)
		}()
	}
	return err
}

/*
The SetObserverOrigin method sets the "origin" value in future Observer messages.

Input
  - ctx: A context to control lifecycle.
  - origin: The value sent in the Observer's "origin" key/value pair.
*/
func (initializerImpl *InitializerImpl) SetObserverOrigin(ctx context.Context, origin string) {
	var err error = nil

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 89
	if initializerImpl.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializerImpl.debug(debugMessageNumber, origin, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializerImpl.getLogger().IsTrace() {
			entryTime := time.Now()
			initializerImpl.traceEntry(80, origin)
			defer func() {
				initializerImpl.traceExit(traceExitMessageNumber, origin, err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJson, err := json.Marshal(initializerImpl)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 81, 1081
			return
		}
		initializerImpl.log(1004, initializerImpl, string(asJson))
	}

	// Set "origin".

	initializerImpl.ObserverOrigin = origin

	senzingSchema := initializerImpl.getSenzingSchema()
	senzingSchema.SetObserverOrigin(ctx, initializerImpl.ObserverOrigin)

	senzingConfig := initializerImpl.getSenzingConfig()
	senzingConfig.SetObserverOrigin(ctx, origin)

	// Notify observers.

	if initializerImpl.observers != nil {
		go func() {
			details := map[string]string{
				"origin": origin,
			}
			notifier.Notify(ctx, initializerImpl.observers, initializerImpl.ObserverOrigin, ComponentId, 8005, err, details)
		}()
	}

}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be removed.
*/
func (initializerImpl *InitializerImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil

	if observer == nil {
		return err
	}

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 79
	if initializerImpl.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializerImpl.debug(debugMessageNumber, observer.GetObserverID(ctx), err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializerImpl.getLogger().IsTrace() {
			entryTime := time.Now()
			initializerImpl.traceEntry(70, observer.GetObserverID(ctx))
			defer func() {
				initializerImpl.traceExit(traceExitMessageNumber, observer.GetObserverID(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJson, err := json.Marshal(initializerImpl)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 71, 1071
			return err
		}
		initializerImpl.log(1005, initializerImpl, string(asJson))
	}

	// Unregister observers in dependencies.

	err = initializerImpl.getSenzingConfig().UnregisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 72, 1072
		return err
	}
	err = initializerImpl.getSenzingSchema().UnregisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 73, 1073
		return err
	}

	// Remove observer from this service.

	if initializerImpl.observers != nil {

		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, initializerImpl.observers, initializerImpl.ObserverOrigin, ComponentId, 8006, err, details)

		err = initializerImpl.observers.UnregisterObserver(ctx, observer)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 74, 1074
			return err
		}

		if !initializerImpl.observers.HasObservers(ctx) {
			initializerImpl.observers = nil
		}
	}
	return err
}
