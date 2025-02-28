package initializer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/senzing-garage/go-databasing/dbhelper"
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

// BasicInitializer is the default implementation of the Initializer interface.
type BasicInitializer struct {
	DataSources           []string `json:"dataSources,omitempty"`
	ObserverOrigin        string   `json:"observerOrigin,omitempty"`
	ObserverURL           string   `json:"observerUrl,omitempty"`
	SenzingInstanceName   string   `json:"senzingInstanceName,omitempty"`
	SenzingLogLevel       string   `json:"senzingLogLevel,omitempty"`
	SenzingSettings       string   `json:"senzingSettings,omitempty"`
	SenzingSettingsFile   string   `json:"senzingSettingsFile,omitempty"`
	SenzingVerboseLogging int64    `json:"senzingVerboseLogging,omitempty"`
	SQLFile               string   `json:"sqlFile,omitempty"`

	logger                 logging.Logging
	observers              subject.Subject
	senzingConfigSingleton senzingconfig.SenzingConfig
	senzingSchemaSingleton senzingschema.SenzingSchema
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var debugOptions = []interface{}{
	&logging.OptionCallerSkip{Value: 5},
}

var traceOptions = []interface{}{
	&logging.OptionCallerSkip{Value: 5},
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
func (initializer *BasicInitializer) Initialize(ctx context.Context) error {
	var err error
	debugMessageNumber := 0
	traceExitMessageNumber := 19

	// Initialize logging.

	logLevel := initializer.SenzingLogLevel
	if logLevel == "" {
		logLevel = "INFO"
	}
	err = initializer.SetLogLevel(ctx, logLevel)
	if err != nil {
		return err
	}

	// Prolog.

	if initializer.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializer.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializer.getLogger().IsTrace() {
			entryTime := time.Now()
			initializer.traceEntry(10)
			defer func() { initializer.traceExit(traceExitMessageNumber, err, time.Since(entryTime)) }()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(initializer)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 11, 1011
			return err
		}
		initializer.log(1000, initializer, string(asJSON))
	}

	// Initialize observing.

	var anObserver observer.Observer
	if len(initializer.ObserverURL) > 0 {
		parsedURL, err := url.Parse(initializer.ObserverURL)
		if err != nil {
			return err
		}
		switch parsedURL.Scheme {
		case "grpc":
			anObserver, err = initializer.createGrpcObserver(ctx, *parsedURL)
			if err != nil {
				traceExitMessageNumber, debugMessageNumber = 18, 1018
				return err
			}
		default:
		}
		err = initializer.registerObserverLocal(ctx, anObserver)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 17, 1017
			return err
		}

		// Notify observers.

		go func() {
			details := map[string]string{
				"observerID": anObserver.GetObserverID(ctx),
			}
			notifier.Notify(ctx, initializer.observers, initializer.ObserverOrigin, ComponentID, 8001, err, details)
		}()
	}

	// Verify database file exists.

	if len(initializer.SQLFile) > 0 {
		_, err = os.Stat(initializer.SQLFile)
		if err != nil {
			initializer.log(3001, initializer.SQLFile)
			traceExitMessageNumber, debugMessageNumber = 21, 1075
			return err
		}
	}

	// Perform initialization for specific databases.

	err = initializer.InitializeSpecificDatabase(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 12, 1012
		return err
	}

	// Create schema in database.

	senzingSchema := initializer.getSenzingSchema()
	err = senzingSchema.SetLogLevel(ctx, logLevel)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 13, 1013
		return err
	}
	err = initializer.registerObserverSenzingSchema(ctx, anObserver)
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

	senzingConfig := initializer.getSenzingConfig()
	err = senzingConfig.SetLogLevel(ctx, logLevel)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 15, 1015
		return err
	}
	err = initializer.registerObserverSenzingConfig(ctx, anObserver)
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

	if initializer.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, initializer.observers, initializer.ObserverOrigin, ComponentID, 8002, err, details)
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
func (initializer *BasicInitializer) InitializeSpecificDatabase(ctx context.Context) error {
	var err error
	var databaseURLs []string

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 49
	if initializer.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializer.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializer.getLogger().IsTrace() {
			entryTime := time.Now()
			initializer.traceEntry(40)
			defer func() { initializer.traceExit(traceExitMessageNumber, err, time.Since(entryTime)) }()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(initializer)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 41, 1041
			return err
		}
		initializer.log(1001, initializer, string(asJSON))
	}

	// Pull values out of SenzingEngineConfigurationJson.

	parser, err := settingsparser.New(initializer.SenzingSettings)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 42, 1042
		return err
	}

	databaseURLs, err = parser.GetDatabaseURLs(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 43, 1043
		return err
	}

	// Process each database.

	for _, databaseURL := range databaseURLs {

		// Parse URL.

		parsedURL, err := dbhelper.ParseDatabaseURL(databaseURL)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 44, 1044
			return err
		}

		// Special handling for each database type.

		switch parsedURL.Scheme {
		case "sqlite3":
			err = initializer.initializeSpecificDatabaseSqlite(ctx, parsedURL)
			if err != nil {
				traceExitMessageNumber, debugMessageNumber = 45, 1045
				return err
			}
		default:
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
func (initializer *BasicInitializer) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if observer == nil {
		return err
	}

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 59
	if initializer.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializer.debug(debugMessageNumber, observer.GetObserverID(ctx), err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializer.getLogger().IsTrace() {
			entryTime := time.Now()
			initializer.traceEntry(50, observer.GetObserverID(ctx))
			defer func() {
				initializer.traceExit(traceExitMessageNumber, observer.GetObserverID(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(initializer)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 51, 1051
			return err
		}
		initializer.log(1002, initializer, string(asJSON))
	}

	// Create empty list of observers.

	if initializer.observers == nil {
		initializer.observers = &subject.SimpleSubject{}
	}

	// Register observer with initializer and dependent services.

	err = initializer.observers.RegisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 52, 1052
		return err
	}
	err = initializer.getSenzingConfig().RegisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 53, 1053
		return err
	}
	err = initializer.getSenzingSchema().RegisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 54, 1054
		return err
	}

	// Notify observers.

	go func() {
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, initializer.observers, initializer.ObserverOrigin, ComponentID, 8003, err, details)
	}()
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (initializer *BasicInitializer) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 69
	if initializer.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializer.debug(debugMessageNumber, logLevelName, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializer.getLogger().IsTrace() {
			entryTime := time.Now()
			initializer.traceEntry(60, logLevelName)
			defer func() { initializer.traceExit(traceExitMessageNumber, logLevelName, err, time.Since(entryTime)) }()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(initializer)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 61, 1061
			return err
		}
		initializer.log(1003, initializer, string(asJSON))
	}

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		traceExitMessageNumber, debugMessageNumber = 62, 1062
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}

	// Set initializer log level.

	err = initializer.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 63, 1063
		return err
	}

	// Set log level for dependent services.

	if initializer.senzingConfigSingleton != nil {
		err = initializer.senzingConfigSingleton.SetLogLevel(ctx, logLevelName)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 64, 1064
			return err
		}
	}
	err = initializer.getSenzingSchema().SetLogLevel(ctx, logLevelName)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 65, 1065
		return err
	}

	// Notify observers.

	if initializer.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, initializer.observers, initializer.ObserverOrigin, ComponentID, 8004, err, details)
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
func (initializer *BasicInitializer) SetObserverOrigin(ctx context.Context, origin string) {
	var err error

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 89
	if initializer.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializer.debug(debugMessageNumber, origin, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializer.getLogger().IsTrace() {
			entryTime := time.Now()
			initializer.traceEntry(80, origin)
			defer func() {
				initializer.traceExit(traceExitMessageNumber, origin, err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(initializer)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 81, 1081
			return
		}
		initializer.log(1004, initializer, string(asJSON))
	}

	// Set "origin".

	initializer.ObserverOrigin = origin

	senzingSchema := initializer.getSenzingSchema()
	senzingSchema.SetObserverOrigin(ctx, initializer.ObserverOrigin)

	senzingConfig := initializer.getSenzingConfig()
	senzingConfig.SetObserverOrigin(ctx, origin)

	// Notify observers.

	if initializer.observers != nil {
		go func() {
			details := map[string]string{
				"origin": origin,
			}
			notifier.Notify(ctx, initializer.observers, initializer.ObserverOrigin, ComponentID, 8005, err, details)
		}()
	}

}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be removed.
*/
func (initializer *BasicInitializer) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if observer == nil {
		return err
	}

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 79
	if initializer.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializer.debug(debugMessageNumber, observer.GetObserverID(ctx), err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializer.getLogger().IsTrace() {
			entryTime := time.Now()
			initializer.traceEntry(70, observer.GetObserverID(ctx))
			defer func() {
				initializer.traceExit(traceExitMessageNumber, observer.GetObserverID(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(initializer)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 71, 1071
			return err
		}
		initializer.log(1005, initializer, string(asJSON))
	}

	// Unregister observers in dependencies.

	err = initializer.getSenzingConfig().UnregisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 72, 1072
		return err
	}
	err = initializer.getSenzingSchema().UnregisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 73, 1073
		return err
	}

	// Remove observer from this service.

	if initializer.observers != nil {

		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, initializer.observers, initializer.ObserverOrigin, ComponentID, 8006, err, details)

		err = initializer.observers.UnregisterObserver(ctx, observer)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 74, 1074
			return err
		}

		if !initializer.observers.HasObservers(ctx) {
			initializer.observers = nil
		}
	}
	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (initializer *BasicInitializer) getLogger() logging.Logging {
	var err error
	if initializer.logger == nil {
		options := []interface{}{
			logging.OptionCallerSkip{Value: 4},
		}
		if len(initializer.SenzingLogLevel) > 0 {
			options = append(options, logging.OptionLogLevel{Value: initializer.SenzingLogLevel})
		}
		initializer.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return initializer.logger
}

// Log message.
func (initializer *BasicInitializer) log(messageNumber int, details ...interface{}) {
	initializer.getLogger().Log(messageNumber, details...)
}

// Debug.
func (initializer *BasicInitializer) debug(messageNumber int, details ...interface{}) {
	details = append(details, debugOptions...)
	initializer.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (initializer *BasicInitializer) traceEntry(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	initializer.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (initializer *BasicInitializer) traceExit(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	initializer.getLogger().Log(messageNumber, details...)
}

// --- Observing --------------------------------------------------------------

func (initializer *BasicInitializer) createGrpcObserver(ctx context.Context, parsedURL url.URL) (observer.Observer, error) {
	_ = ctx
	var err error

	var result observer.Observer

	port := DefaultGrpcObserverPort
	if len(parsedURL.Port()) > 0 {
		port = parsedURL.Port()
	}
	target := fmt.Sprintf("%s:%s", parsedURL.Hostname(), port)

	// TODO: Allow specification of options from ObserverUrl/parsedUrl
	grpcOptions := grpc.WithTransportCredentials(insecure.NewCredentials())

	grpcConnection, err := grpc.NewClient(target, grpcOptions)
	if err != nil {
		return result, err
	}
	result = &observer.GrpcObserver{
		GrpcClient: observerpb.NewObserverClient(grpcConnection),
		ID:         "init-database",
	}
	return result, err
}

func (initializer *BasicInitializer) registerObserverLocal(ctx context.Context, observer observer.Observer) error {
	if initializer.observers == nil {
		initializer.observers = &subject.SimpleSubject{}
	}
	return initializer.observers.RegisterObserver(ctx, observer)
}

func (initializer *BasicInitializer) registerObserverSenzingConfig(ctx context.Context, observer observer.Observer) error {
	initializer.getSenzingConfig().SetObserverOrigin(ctx, initializer.ObserverOrigin)
	return initializer.getSenzingConfig().RegisterObserver(ctx, observer)
}

func (initializer *BasicInitializer) registerObserverSenzingSchema(ctx context.Context, observer observer.Observer) error {
	initializer.getSenzingSchema().SetObserverOrigin(ctx, initializer.ObserverOrigin)
	return initializer.getSenzingSchema().RegisterObserver(ctx, observer)
}

// --- Dependent services -----------------------------------------------------

func (initializer *BasicInitializer) getSenzingConfig() senzingconfig.SenzingConfig {
	if initializer.senzingConfigSingleton == nil {
		initializer.senzingConfigSingleton = &senzingconfig.BasicSenzingConfig{
			DataSources:           initializer.DataSources,
			SenzingSettingsFile:   initializer.SenzingSettingsFile,
			SenzingSettings:       initializer.SenzingSettings,
			SenzingInstanceName:   initializer.SenzingInstanceName,
			SenzingVerboseLogging: initializer.SenzingVerboseLogging,
		}
	}
	return initializer.senzingConfigSingleton
}

func (initializer *BasicInitializer) getSenzingSchema() senzingschema.SenzingSchema {
	if initializer.senzingSchemaSingleton == nil {
		initializer.senzingSchemaSingleton = &senzingschema.BasicSenzingSchema{
			SenzingSettings: initializer.SenzingSettings,
			SQLFile:         initializer.SQLFile,
		}
	}
	return initializer.senzingSchemaSingleton
}

// --- Specific database processing -------------------------------------------

func (initializer *BasicInitializer) initializeSpecificDatabaseSqlite(ctx context.Context, parsedURL *url.URL) error {
	var err error

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 109
	if initializer.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				initializer.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if initializer.getLogger().IsTrace() {
			entryTime := time.Now()
			initializer.traceEntry(100, parsedURL)
			defer func() { initializer.traceExit(traceExitMessageNumber, parsedURL, err, time.Since(entryTime)) }()
		}
	}

	// If in-memory database, do not create a file.

	queryParameters := parsedURL.Query()
	if (queryParameters.Get("mode") == "memory") && (queryParameters.Get("cache") == "shared") {
		return err // Nothing to do for in-memory database.
	}

	// If file exists, no more to do.

	filename := parsedURL.Path
	filename = filepath.Clean(filename) // See https://securego.io/docs/rules/g304.html
	filename = cleanFilename(filename)
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
	initializer.log(2001, filename)

	// Notify observers.

	if initializer.observers != nil {
		go func() {
			details := map[string]string{
				"sqliteFile": filename,
			}
			notifier.Notify(ctx, initializer.observers, initializer.ObserverOrigin, ComponentID, 8010, err, details)
		}()
	}
	return err
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
