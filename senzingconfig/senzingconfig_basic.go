package senzingconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing-garage/go-sdk-abstract-factory/szfactorycreator"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicSenzingConfig is the default implementation of the SenzingConfig interface.
type BasicSenzingConfig struct {
	DataSources           []string          `json:"dataSources,omitempty"`
	GrpcDialOptions       []grpc.DialOption `json:"grpcDialOptions,omitempty"`
	GrpcTarget            string            `json:"grpcTarget,omitempty"`
	SenzingInstanceName   string            `json:"senzingInstanceName,omitempty"`
	SenzingSettings       string            `json:"senzingSettings,omitempty"`
	SenzingConfigJSONFile string            `json:"senzingConfigJsonFile,omitempty"`
	SenzingVerboseLogging int64             `json:"senzingVerboseLogging,omitempty"`

	isTrace                    bool
	logger                     logging.Logging
	logLevel                   string
	observerOrigin             string
	observers                  subject.Subject
	szAbstractFactorySingleton senzing.SzAbstractFactory
	szAbstractFactorySyncOnce  sync.Once
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var debugOptions = []interface{}{
	&logging.OptionCallerSkip{Value: OptionCallerSkip5},
}

var traceOptions = []interface{}{
	&logging.OptionCallerSkip{Value: OptionCallerSkip5},
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The InitializeSenzing method adds the Senzing default configuration to databases.

Input
  - ctx: A context to control lifecycle.
*/
func (senzingConfig *BasicSenzingConfig) InitializeSenzing(ctx context.Context) error {
	var err error

	var configID int64

	entryTime := time.Now()

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 29

	if senzingConfig.getLogger().IsDebug() {
		// If DEBUG, log error exit.
		defer func() {
			if debugMessageNumber > 0 {
				senzingConfig.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingConfig.getLogger().IsTrace() {
			senzingConfig.traceEntry(10, configID)

			defer func() { senzingConfig.traceExit(traceExitMessageNumber, configID, err, time.Since(entryTime)) }()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingConfig)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 11, 1011

			return wraperror.Errorf(err, "json.Marshal: %v", senzingConfig)
		}

		senzingConfig.log(1001, senzingConfig, string(asJSON))
	}

	// Create Senzing objects.

	szAbstractFactory := senzingConfig.getAbstractFactory(ctx)

	defer func() { szAbstractFactory.Close(ctx) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		return wraperror.Errorf(err, "CreateConfigManager")
	}

	defer func() { _ = szConfigManager.Destroy(ctx) }()

	// If a Senzing configuration file is specified, use it.

	if len(senzingConfig.SenzingConfigJSONFile) > 0 {
		configDefinition, err1 := fileToString(ctx, senzingConfig.SenzingConfigJSONFile)
		if err1 != nil {
			traceExitMessageNumber, debugMessageNumber = 99, 1999

			return wraperror.Errorf(err1, "fileToString: %s", senzingConfig.SenzingConfigJSONFile)
		}

		szConfig, err2 := szConfigManager.CreateConfigFromString(ctx, configDefinition)
		if err2 != nil {
			traceExitMessageNumber, debugMessageNumber = 99, 1999

			return wraperror.Errorf(err2, "CreateConfigFromString: %s", configDefinition)
		}

		configID, err = senzingConfig.makeDefaultConfig(ctx, szAbstractFactory, szConfig)
		senzingConfig.log(2999, configID)

		traceExitMessageNumber, debugMessageNumber = 99, 999

		return wraperror.Errorf(err, "makeDefaultConfig")
	}

	// Determine if configuration already exists. If so, return.

	configID, err = szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 13, 1013

		return wraperror.Errorf(err, "GetDefaultConfigID")
	}

	if configID != 0 {
		if senzingConfig.observers != nil {
			go func() {
				details := map[string]string{}
				notifier.Notify(
					ctx,
					senzingConfig.observers,
					senzingConfig.observerOrigin,
					ComponentID,
					8001,
					err,
					details,
				)
			}()
		}

		if len(senzingConfig.DataSources) > 0 {
			szConfig, err2 := szConfigManager.CreateConfigFromConfigID(ctx, configID)
			if err2 != nil {
				traceExitMessageNumber, debugMessageNumber = 99, 1999

				return wraperror.Errorf(err2, "CreateConfigFromConfigID: %d", configID)
			}

			configID, err = senzingConfig.makeDefaultConfig(ctx, szAbstractFactory, szConfig)
		}

		senzingConfig.log(2002, configID)

		traceExitMessageNumber, debugMessageNumber = 14, 0 // debugMessageNumber=0 because it's not an error.

		return wraperror.Errorf(err, "ConfigID: %d", configID)
	}

	// If no configuration file specified, install the template.

	if len(senzingConfig.SenzingConfigJSONFile) == 0 {
		szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 99, 1999

			return wraperror.Errorf(err, "CreateConfigFromTemplate")
		}

		configID, err = senzingConfig.makeDefaultConfig(ctx, szAbstractFactory, szConfig)
		senzingConfig.log(2999, configID)

		traceExitMessageNumber, debugMessageNumber = 999, 999

		return wraperror.Errorf(err, "makeDefaultConfig")
	}

	// Notify observers.

	senzingConfig.log(2003, configID)

	if senzingConfig.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentID, 8002, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (senzingConfig *BasicSenzingConfig) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if observer == nil {
		return err
	}

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 39

	if senzingConfig.getLogger().IsDebug() {
		// If DEBUG, log error exit.
		defer func() {
			if debugMessageNumber > 0 {
				senzingConfig.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingConfig.getLogger().IsTrace() {
			entryTime := time.Now()

			senzingConfig.traceEntry(30, observer.GetObserverID(ctx))

			defer func() {
				senzingConfig.traceExit(traceExitMessageNumber, observer.GetObserverID(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingConfig)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 31, 1031

			return wraperror.Errorf(err, "json.Marshal: %v", senzingConfig)
		}

		senzingConfig.log(1002, senzingConfig, string(asJSON))
	}

	// Create empty list of observers.

	if senzingConfig.observers == nil {
		senzingConfig.observers = &subject.SimpleSubject{}
	}

	// Notify observers.

	go func() {
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentID, 8003, err, details)
	}()

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (senzingConfig *BasicSenzingConfig) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 49

	if senzingConfig.getLogger().IsDebug() {
		// If DEBUG, log error exit.
		defer func() {
			if debugMessageNumber > 0 {
				senzingConfig.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingConfig.getLogger().IsTrace() {
			entryTime := time.Now()

			senzingConfig.traceEntry(40, logLevelName)

			defer func() {
				senzingConfig.traceExit(traceExitMessageNumber, logLevelName, err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingConfig)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 41, 1041

			return wraperror.Errorf(err, "json.Marshal: %v", senzingConfig)
		}

		senzingConfig.log(1003, senzingConfig, string(asJSON))
	}

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		traceExitMessageNumber, debugMessageNumber = 42, 1042

		return wraperror.Errorf(errForPackage, "invalid error level: %s", logLevelName)
	}

	// Set senzingConfig log level.

	senzingConfig.logLevel = logLevelName

	err = senzingConfig.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 43, 1043

		return wraperror.Errorf(err, "SetLogLevel: %s", logLevelName)
	}

	senzingConfig.isTrace = (logLevelName == logging.LevelTraceName)

	// Notify observers.

	if senzingConfig.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentID, 8004, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
The SetObserverOrigin method sets the "origin" value in future Observer messages.

Input
  - ctx: A context to control lifecycle.
  - origin: The value sent in the Observer's "origin" key/value pair.
*/
func (senzingConfig *BasicSenzingConfig) SetObserverOrigin(ctx context.Context, origin string) {
	var err error

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 69

	if senzingConfig.getLogger().IsDebug() {
		// If DEBUG, log error exit.
		defer func() {
			if debugMessageNumber > 0 {
				senzingConfig.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingConfig.getLogger().IsTrace() {
			entryTime := time.Now()

			senzingConfig.traceEntry(60, origin)

			defer func() {
				senzingConfig.traceExit(traceExitMessageNumber, origin, err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingConfig)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 61, 1061

			return
		}

		senzingConfig.log(1004, senzingConfig, string(asJSON))
	}

	// Notify observers.

	if senzingConfig.observers != nil {
		go func() {
			details := map[string]string{
				"origin": origin,
			}
			notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentID, 8005, err, details)
		}()
	}
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be removed.
*/
func (senzingConfig *BasicSenzingConfig) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if observer == nil {
		return err
	}

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 59

	if senzingConfig.getLogger().IsDebug() {
		// If DEBUG, log error exit.
		defer func() {
			if debugMessageNumber > 0 {
				senzingConfig.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingConfig.getLogger().IsTrace() {
			entryTime := time.Now()

			senzingConfig.traceEntry(50, observer.GetObserverID(ctx))

			defer func() {
				senzingConfig.traceExit(traceExitMessageNumber, observer.GetObserverID(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingConfig)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 51, 1051

			return wraperror.Errorf(err, "json.Marshal: %v", senzingConfig)
		}

		senzingConfig.log(1005, senzingConfig, string(asJSON))
	}

	// Remove observer from this service.

	if senzingConfig.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentID, 8006, err, details)

		err = senzingConfig.observers.UnregisterObserver(ctx, observer)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 54, 1054

			return wraperror.Errorf(err, "UnregisterObserver")
		}

		if !senzingConfig.observers.HasObservers(ctx) {
			senzingConfig.observers = nil
		}
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (senzingConfig *BasicSenzingConfig) getLogger() logging.Logging {
	var err error

	if senzingConfig.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: OptionCallerSkip4},
		}

		senzingConfig.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}

	return senzingConfig.logger
}

// Log message.
func (senzingConfig *BasicSenzingConfig) log(messageNumber int, details ...interface{}) {
	senzingConfig.getLogger().Log(messageNumber, details...)
}

// Debug.
func (senzingConfig *BasicSenzingConfig) debug(messageNumber int, details ...interface{}) {
	details = append(details, debugOptions...)
	senzingConfig.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (senzingConfig *BasicSenzingConfig) traceEntry(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	senzingConfig.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (senzingConfig *BasicSenzingConfig) traceExit(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	senzingConfig.getLogger().Log(messageNumber, details...)
}

// --- Dependent services -----------------------------------------------------

// Create an abstract factory singleton and return it.
func (senzingConfig *BasicSenzingConfig) getAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	var err error

	_ = ctx

	senzingConfig.szAbstractFactorySyncOnce.Do(func() {
		if len(senzingConfig.GrpcTarget) == 0 {
			senzingConfig.szAbstractFactorySingleton, err = senzingConfig.buildSzAbstractFactory()
			if err != nil {
				panic(err)
			}
		} else {
			grpcConnection, err := grpc.NewClient(senzingConfig.GrpcTarget, senzingConfig.GrpcDialOptions...)
			if err != nil {
				panic(err)
			}

			senzingConfig.szAbstractFactorySingleton, err = szfactorycreator.CreateGrpcAbstractFactory(grpcConnection)
			if err != nil {
				panic(err)
			}
		}
	})

	return senzingConfig.szAbstractFactorySingleton
}

// --- Misc -------------------------------------------------------------------

func (senzingConfig *BasicSenzingConfig) buildSzAbstractFactory() (senzing.SzAbstractFactory, error) {
	var (
		err    error
		result senzing.SzAbstractFactory
	)

	senzingSettings := senzingConfig.SenzingSettings

	senzingInstanceName := senzingConfig.SenzingInstanceName
	if len(senzingInstanceName) == 0 {
		senzingInstanceName = fmt.Sprintf("senzing init-database at %s", time.Now())
	}

	result, err = szfactorycreator.CreateCoreAbstractFactory(
		senzingInstanceName,
		senzingSettings,
		senzingConfig.SenzingVerboseLogging,
		senzing.SzInitializeWithDefaultConfiguration,
	)

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

func (senzingConfig *BasicSenzingConfig) makeDefaultConfig(
	ctx context.Context,
	szAbstractFactory senzing.SzAbstractFactory,
	szConfig senzing.SzConfig,
) (int64, error) {
	var (
		err    error
		result int64
	)

	for _, datasource := range senzingConfig.DataSources {
		_, err = szConfig.RegisterDataSource(ctx, datasource)
		if err != nil {
			return result, wraperror.Errorf(err, "RegisterDataSource: %s", datasource)
		}

		senzingConfig.log(2001, datasource)
	}

	configComment := fmt.Sprintf(
		"Created by init-database at %s with datasources: %s ",
		time.Now().Format(time.RFC3339),
		strings.Join(senzingConfig.DataSources, " "),
	)

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		return result, wraperror.Errorf(err, "Export")
	}

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		return result, wraperror.Errorf(err, "CreateConfigManager")
	}

	defer func() { _ = szConfigManager.Destroy(ctx) }()

	result, err = szConfigManager.SetDefaultConfig(ctx, configDefinition, configComment)
	if err != nil {
		return result, wraperror.Errorf(err, "SetDefaultConfig")
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func fileToString(ctx context.Context, filePath string) (string, error) {
	_ = ctx
	content, err := os.ReadFile(filepath.Clean(filePath))

	return string(content), err
}

// func assertNoError(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
