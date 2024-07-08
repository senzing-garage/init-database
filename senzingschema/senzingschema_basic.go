package senzingschema

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/senzing-garage/go-databasing/connector"
	"github.com/senzing-garage/go-databasing/sqlexecutor"
	"github.com/senzing-garage/go-helpers/settingsparser"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicSenzingSchema is the default implementation of the SenzingSchema interface.
type BasicSenzingSchema struct {
	SenzingSettings string `json:"senzingSettings,omitempty"`
	SQLFile         string `json:"sqlFile,omitempty"`

	logger         logging.Logging
	logLevelName   string
	observerOrigin string
	observers      subject.Subject
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
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (senzingSchema *BasicSenzingSchema) getLogger() logging.Logging {
	var err error
	if senzingSchema.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		senzingSchema.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return senzingSchema.logger
}

// Log message.
func (senzingSchema *BasicSenzingSchema) log(messageNumber int, details ...interface{}) {
	senzingSchema.getLogger().Log(messageNumber, details...)
}

// Debug.
func (senzingSchema *BasicSenzingSchema) debug(messageNumber int, details ...interface{}) {
	details = append(details, debugOptions...)
	senzingSchema.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (senzingSchema *BasicSenzingSchema) traceEntry(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	senzingSchema.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (senzingSchema *BasicSenzingSchema) traceExit(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	senzingSchema.getLogger().Log(messageNumber, details...)
}

// --- Misc -------------------------------------------------------------------

// Given a database URL, detemine the correct SQL file and send the statements to the database.
func (senzingSchema *BasicSenzingSchema) processDatabase(ctx context.Context, resourcePath string, databaseURL string) error {
	var err error

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 109
	if senzingSchema.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				senzingSchema.debug(debugMessageNumber, resourcePath, databaseURL, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingSchema.getLogger().IsTrace() {
			entryTime := time.Now()
			senzingSchema.traceEntry(100, resourcePath, databaseURL)
			defer func() {
				senzingSchema.traceExit(traceExitMessageNumber, resourcePath, databaseURL, err, time.Since(entryTime))
			}()
		}
	}

	// Determine which SQL file to process.

	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		if strings.HasPrefix(databaseURL, "postgresql") {
			index := strings.LastIndex(databaseURL, ":")
			newDatabaseURL := databaseURL[:index] + "/" + databaseURL[index+1:]
			parsedURL, err = url.Parse(newDatabaseURL)
		}
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 101, 1101
			return err
		}
	}

	if len(senzingSchema.SQLFile) == 0 {
		switch parsedURL.Scheme {
		case "sqlite3":
			senzingSchema.SQLFile = resourcePath + "/schema/g2core-schema-sqlite-create.sql"
		case "postgresql":
			senzingSchema.SQLFile = resourcePath + "/schema/g2core-schema-postgresql-create.sql"
		case "mysql":
			senzingSchema.SQLFile = resourcePath + "/schema/g2core-schema-mysql-create.sql"
		case "mssql":
			senzingSchema.SQLFile = resourcePath + "/schema/g2core-schema-mssql-create.sql"
		default:
			return fmt.Errorf("unknown database scheme: %s", parsedURL.Scheme)
		}
	}

	// Connect to the database.

	databaseConnector, err := connector.NewConnector(ctx, databaseURL)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 102, 1102
		return err
	}

	// Create sqlExecutor to process file of SQL.

	sqlExecutor := &sqlexecutor.BasicSQLExecutor{
		DatabaseConnector: databaseConnector,
	}
	err = sqlExecutor.SetLogLevel(ctx, senzingSchema.logLevelName)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 103, 1103
		return err
	}

	// Add observers to sqlExecutor.

	if senzingSchema.observers != nil {
		for _, observer := range senzingSchema.observers.GetObservers(ctx) {
			err = sqlExecutor.RegisterObserver(ctx, observer)
			if err != nil {
				traceExitMessageNumber, debugMessageNumber = 104, 1104
				return err
			}
		}
	}

	// TODO: add following when it becomes available.
	// sqlExecutor.SetObserverOrigin(ctx, senzingSchema.observerOrigin)

	// Process file of SQL

	err = sqlExecutor.ProcessFileName(ctx, senzingSchema.SQLFile)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 105, 1105
		return err
	}
	senzingSchema.log(2001, senzingSchema.SQLFile, parsedURL.Redacted())
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
func (senzingSchema *BasicSenzingSchema) InitializeSenzing(ctx context.Context) error {
	var err error

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 19
	if senzingSchema.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				senzingSchema.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingSchema.getLogger().IsTrace() {
			entryTime := time.Now()
			senzingSchema.traceEntry(10)
			defer func() { senzingSchema.traceExit(traceExitMessageNumber, err, time.Since(entryTime)) }()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingSchema)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 11, 1011
			return err
		}
		senzingSchema.log(1001, senzingSchema, string(asJSON))
	}

	// Pull values out of SenzingEngineConfigurationJson.

	parser, err := settingsparser.New(senzingSchema.SenzingSettings)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 12, 1012
		return err
	}
	resourcePath, err := parser.GetResourcePath(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 13, 1013
		return err
	}
	databaseURLs, err := parser.GetDatabaseURLs(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 14, 1014
		return err
	}

	// Process each database.

	for _, databaseURL := range databaseURLs {
		err = senzingSchema.processDatabase(ctx, resourcePath, databaseURL)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 15, 1015
			return err
		}
	}

	// Notify observers.

	if senzingSchema.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, senzingSchema.observers, senzingSchema.observerOrigin, ComponentID, 8001, err, details)
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
func (senzingSchema *BasicSenzingSchema) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if observer == nil {
		return err
	}

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 29
	if senzingSchema.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				senzingSchema.debug(debugMessageNumber, observer.GetObserverID(ctx), err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingSchema.getLogger().IsTrace() {
			entryTime := time.Now()
			senzingSchema.traceEntry(20, observer.GetObserverID(ctx))
			defer func() {
				senzingSchema.traceExit(traceExitMessageNumber, observer.GetObserverID(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingSchema)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 21, 1021
			return err
		}
		senzingSchema.log(1002, senzingSchema, string(asJSON))
	}

	// Create empty list of observers.

	if senzingSchema.observers == nil {
		senzingSchema.observers = &subject.SimpleSubject{}
	}

	// Register observer with senzingSchema.

	err = senzingSchema.observers.RegisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 22, 1022
		return err
	}

	// Notify observers.

	go func() {
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, senzingSchema.observers, senzingSchema.observerOrigin, ComponentID, 8002, err, details)
	}()

	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (senzingSchema *BasicSenzingSchema) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 39
	if senzingSchema.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				senzingSchema.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingSchema.getLogger().IsTrace() {
			entryTime := time.Now()
			senzingSchema.traceEntry(30, logLevelName)
			defer func() {
				senzingSchema.traceExit(traceExitMessageNumber, logLevelName, err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingSchema)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 31, 1031
			return err
		}
		senzingSchema.log(1003, senzingSchema, string(asJSON))
	}

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		traceExitMessageNumber, debugMessageNumber = 32, 1032
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}

	// Set senzingSchema log level.

	senzingSchema.logLevelName = logLevelName
	err = senzingSchema.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 33, 1033
		return err
	}

	// Notify observers.

	if senzingSchema.observers != nil { // Performance optimization.
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, senzingSchema.observers, senzingSchema.observerOrigin, ComponentID, 8003, err, details)
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
func (senzingSchema *BasicSenzingSchema) SetObserverOrigin(ctx context.Context, origin string) {
	var err error

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 59
	if senzingSchema.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				senzingSchema.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingSchema.getLogger().IsTrace() {
			entryTime := time.Now()
			senzingSchema.traceEntry(50, origin)
			defer func() {
				senzingSchema.traceExit(traceExitMessageNumber, origin, err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingSchema)
		if err != nil {
			debugMessageNumber = 1051
			traceExitMessageNumber = 51
			traceExitMessageNumber, debugMessageNumber = 51, 1051
			return
		}
		senzingSchema.log(1004, senzingSchema, string(asJSON))
	}

	// Set origin in dependent services.

	senzingSchema.observerOrigin = origin

	// Notify observers.

	if senzingSchema.observers != nil {
		go func() {
			details := map[string]string{
				"origin": origin,
			}
			notifier.Notify(ctx, senzingSchema.observers, senzingSchema.observerOrigin, ComponentID, 8004, err, details)
		}()
	}

}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be removed.
*/
func (senzingSchema *BasicSenzingSchema) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if observer == nil {
		return err
	}

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 49
	if senzingSchema.getLogger().IsDebug() {

		// If DEBUG, log error exit.

		defer func() {
			if debugMessageNumber > 0 {
				senzingSchema.debug(debugMessageNumber, observer.GetObserverID(ctx), err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingSchema.getLogger().IsTrace() {
			entryTime := time.Now()
			senzingSchema.traceEntry(40, observer.GetObserverID(ctx))
			defer func() {
				senzingSchema.traceExit(traceExitMessageNumber, observer.GetObserverID(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingSchema)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 41, 1041
			return err
		}
		senzingSchema.log(1005, senzingSchema, string(asJSON))
	}

	// Remove observer from this service.

	if senzingSchema.observers != nil {

		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, senzingSchema.observers, senzingSchema.observerOrigin, ComponentID, 8005, err, details)

		err = senzingSchema.observers.UnregisterObserver(ctx, observer)
		if err != nil {
			debugMessageNumber = 1042
			traceExitMessageNumber = 42
			traceExitMessageNumber, debugMessageNumber = 42, 1042
			return err
		}

		if !senzingSchema.observers.HasObservers(ctx) {
			senzingSchema.observers = nil
		}
	}

	return err
}
