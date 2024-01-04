package senzingconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/senzing-garage/go-common/engineconfigurationjsonparser"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing/go-sdk-abstract-factory/factory"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// SenzingConfigImpl is the default implementation of the SenzingConfig interface.
type SenzingConfigImpl struct {
	DataSources                    []string
	g2configmgrSingleton           g2api.G2configmgr
	g2configmgrSyncOnce            sync.Once
	g2configSingleton              g2api.G2config
	g2configSyncOnce               sync.Once
	g2factorySingleton             factory.SdkAbstractFactory
	g2factorySyncOnce              sync.Once
	isTrace                        bool
	logger                         logging.LoggingInterface
	logLevel                       string
	observerOrigin                 string
	observers                      subject.Subject
	SenzingEngineConfigurationFile string
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int64
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

var defaultModuleName string = "init-database"

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (senzingConfig *SenzingConfigImpl) getLogger() logging.LoggingInterface {
	var err error = nil
	if senzingConfig.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		senzingConfig.logger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return senzingConfig.logger
}

// Log message.
func (senzingConfig *SenzingConfigImpl) log(messageNumber int, details ...interface{}) {
	senzingConfig.getLogger().Log(messageNumber, details...)
}

// Debug.
func (senzingConfig *SenzingConfigImpl) debug(messageNumber int, details ...interface{}) {
	details = append(details, debugOptions...)
	senzingConfig.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (senzingConfig *SenzingConfigImpl) traceEntry(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	senzingConfig.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (senzingConfig *SenzingConfigImpl) traceExit(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	senzingConfig.getLogger().Log(messageNumber, details...)
}

// --- Dependent services -----------------------------------------------------

// Create an abstract factory singleton and return it.
func (senzingConfig *SenzingConfigImpl) getG2Factory(ctx context.Context) factory.SdkAbstractFactory {
	senzingConfig.g2factorySyncOnce.Do(func() {
		senzingConfig.g2factorySingleton = &factory.SdkAbstractFactoryImpl{}
	})
	return senzingConfig.g2factorySingleton
}

// Create a G2Config singleton and return it.
func (senzingConfig *SenzingConfigImpl) getG2config(ctx context.Context) (g2api.G2config, error) {
	var err error = nil
	senzingConfig.g2configSyncOnce.Do(func() {
		senzingConfig.g2configSingleton, err = senzingConfig.getG2Factory(ctx).GetG2config(ctx)
		if err != nil {
			return
		}
		if senzingConfig.g2configSingleton.GetSdkId(ctx) == "base" {
			moduleName := senzingConfig.SenzingModuleName
			if len(moduleName) == 0 {
				moduleName = defaultModuleName
			}
			err = senzingConfig.g2configSingleton.Init(ctx, moduleName, senzingConfig.SenzingEngineConfigurationJson, senzingConfig.SenzingVerboseLogging)
		}
	})
	return senzingConfig.g2configSingleton, err
}

// Create a G2Configmgr singleton and return it.
func (senzingConfig *SenzingConfigImpl) getG2configmgr(ctx context.Context) (g2api.G2configmgr, error) {
	var err error = nil
	senzingConfig.g2configmgrSyncOnce.Do(func() {
		senzingConfig.g2configmgrSingleton, err = senzingConfig.getG2Factory(ctx).GetG2configmgr(ctx)
		if err != nil {
			return
		}
		if senzingConfig.g2configmgrSingleton.GetSdkId(ctx) == "base" {
			moduleName := senzingConfig.SenzingModuleName
			if len(moduleName) == 0 {
				moduleName = defaultModuleName
			}
			err = senzingConfig.g2configmgrSingleton.Init(ctx, moduleName, senzingConfig.SenzingEngineConfigurationJson, senzingConfig.SenzingVerboseLogging)
		}
	})
	return senzingConfig.g2configmgrSingleton, err
}

// Get dependent services: G2config, G2configmgr
func (senzingConfig *SenzingConfigImpl) getDependentServices(ctx context.Context) (g2api.G2config, g2api.G2configmgr, error) {
	g2Config, err := senzingConfig.getG2config(ctx)
	if err != nil {
		return nil, nil, err
	}
	g2Configmgr, err := senzingConfig.getG2configmgr(ctx)
	if err != nil {
		return g2Config, nil, err
	}
	return g2Config, g2Configmgr, err
}

// --- Misc -------------------------------------------------------------------

// Add datasources to Senzing configuration.
func (senzingConfig *SenzingConfigImpl) addDatasources(ctx context.Context, g2Config g2api.G2config, configHandle uintptr) error {
	var err error = nil
	for _, datasource := range senzingConfig.DataSources {
		inputJson := `{"DSRC_CODE": "` + datasource + `"}`
		_, err = g2Config.AddDataSource(ctx, configHandle, inputJson)
		if err != nil {
			return err
		}
		senzingConfig.log(2001, datasource)
	}
	return err
}

func (senzingConfig *SenzingConfigImpl) copyFile(sourceFilename string, targetFilename string) error {
	sourceFilename = filepath.Clean(sourceFilename)
	sourceFile, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			senzingConfig.log(9999, sourceFilename, err)
		}
	}()
	targetFilename = filepath.Clean(targetFilename)
	targetFile, err := os.Create(targetFilename)
	if err != nil {
		return err
	}
	defer func() {
		if err := targetFile.Close(); err != nil {
			senzingConfig.log(9999, targetFilename, err)
		}
	}()
	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return err
	}
	senzingConfig.log(2004, sourceFilename, targetFilename)
	return err
}

func (senzingConfig *SenzingConfigImpl) filesAreEqual(sourceFilename string, targetFilename string) bool {
	var (
		chunkSize        int  = 64000
		shortCircuitExit bool = false
	)

	// If file sizes differ, then files differ.

	sourceStat, err := os.Stat(sourceFilename)
	if err != nil {
		return false
	}
	targetStat, err := os.Stat(targetFilename)
	if err != nil {
		return false
	}

	if sourceStat.Size() != targetStat.Size() {
		return false
	}

	// Final check: If file contents differ, then files differ.

	sourceFilename = filepath.Clean(sourceFilename)
	sourceFile, err := os.Open(sourceFilename)
	if err != nil {
		shortCircuitExit = true
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			senzingConfig.log(9999, sourceFilename, err)
		}
	}()

	targetFilename = filepath.Clean(targetFilename)
	targetFile, err := os.Open(targetFilename)
	if err != nil {
		shortCircuitExit = true
	}
	defer func() {
		if err := targetFile.Close(); err != nil {
			senzingConfig.log(9999, targetFilename, err)
		}
	}()

	if shortCircuitExit {
		return false
	}

	for {
		sourceBytes := make([]byte, chunkSize)
		_, sourceError := sourceFile.Read(sourceBytes)

		targetBytes := make([]byte, chunkSize)
		_, targetError := targetFile.Read(targetBytes)

		if sourceError != nil || targetError != nil {
			if sourceError == io.EOF && targetError == io.EOF {
				return true
			} else if sourceError == io.EOF || targetError == io.EOF {
				return false
			} else {
				senzingConfig.log(4001, sourceFilename, targetFilename, sourceError, targetError)
			}
		}
		if !bytes.Equal(sourceBytes, targetBytes) {
			return false
		}
	}
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The InitializeSenzing method adds the Senzing default configuration to databases.

Input
  - ctx: A context to control lifecycle.
*/
func (senzingConfig *SenzingConfigImpl) InitializeSenzing(ctx context.Context) error {
	var err error = nil
	var configID int64 = 0
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

		asJson, err := json.Marshal(senzingConfig)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 11, 1011
			return err
		}
		senzingConfig.log(1001, senzingConfig, string(asJson))
	}

	// Create Senzing objects.

	g2Config, g2Configmgr, err := senzingConfig.getDependentServices(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 12, 1012
		return err
	}

	// Determine if configuration already exists. If so, return.

	configID, err = g2Configmgr.GetDefaultConfigID(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 13, 1013
		return err
	}
	if configID != 0 {
		if senzingConfig.observers != nil {
			go func() {
				details := map[string]string{}
				notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentId, 8001, err, details)
			}()
		}
		senzingConfig.log(2002, configID)
		traceExitMessageNumber, debugMessageNumber = 14, 0 // debugMessageNumber=0 because it's not an error.
		return err
	}

	// If engine configuration file specified, swap it in.

	if len(senzingConfig.SenzingEngineConfigurationFile) > 0 {
		parsedJson, err := engineconfigurationjsonparser.New(senzingConfig.SenzingEngineConfigurationJson)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 20, 1020
			return err
		}
		resourcePath, err := parsedJson.GetResourcePath(ctx)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 21, 1021
			return err
		}

		// Compare file names.

		sourceFilename := senzingConfig.SenzingEngineConfigurationFile
		targetFilename := fmt.Sprintf("%s/templates/g2config.json", resourcePath)
		if sourceFilename != targetFilename {

			// Verify source file exists.

			_, err := os.Stat(sourceFilename)
			if err != nil {
				senzingConfig.log(5001, sourceFilename, err)
				traceExitMessageNumber, debugMessageNumber = 22, 1022
				return err
			}

			// Determine if target file needs to be replaced.

			if senzingConfig.filesAreEqual(sourceFilename, targetFilename) {
				senzingConfig.log(2005, sourceFilename, targetFilename)
			} else {

				// If target file exists, back it up.

				_, err = os.Stat(targetFilename)
				if err == nil {
					backupFilename := fmt.Sprintf("%s.%d", targetFilename, time.Now().Unix())
					err = senzingConfig.copyFile(targetFilename, backupFilename)
					if err != nil {
						senzingConfig.log(5002, targetFilename, backupFilename, err)
						traceExitMessageNumber, debugMessageNumber = 23, 1023
						return err
					}
				}

				// Copy source file to target to "fake out" Senzing's G2Engine.Create().

				err = senzingConfig.copyFile(sourceFilename, targetFilename)
				if err != nil {
					senzingConfig.log(5003, sourceFilename, targetFilename, err)
					traceExitMessageNumber, debugMessageNumber = 24, 1024
					return err
				}
			}
		}
	}

	// Create a fresh Senzing configuration.

	configHandle, err := g2Config.Create(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 15, 1015
		return err
	}

	// If requested, add DataSources to fresh Senzing configuration.

	if len(senzingConfig.DataSources) > 0 {
		err = senzingConfig.addDatasources(ctx, g2Config, configHandle)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 16, 1016
			return err
		}
	}

	// Create a JSON string from the in-memory configuration.

	configStr, err := g2Config.Save(ctx, configHandle)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 17, 1017
		return err
	}

	// Persist the Senzing configuration to the Senzing repository and set as default configuration.

	configComments := fmt.Sprintf("Created by %s at %s", defaultModuleName, entryTime.Format(time.RFC3339Nano))
	configID, err = g2Configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 18, 1018
		return err
	}
	err = g2Configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 19, 1019
		return err
	}

	// Notify observers.

	senzingConfig.log(2003, configID, configComments)
	if senzingConfig.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentId, 8002, err, details)
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
func (senzingConfig *SenzingConfigImpl) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil

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
			senzingConfig.traceEntry(30, observer.GetObserverId(ctx))
			defer func() {
				senzingConfig.traceExit(traceExitMessageNumber, observer.GetObserverId(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJson, err := json.Marshal(senzingConfig)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 31, 1031
			return err
		}
		senzingConfig.log(1002, senzingConfig, string(asJson))
	}

	// Create empty list of observers.

	if senzingConfig.observers == nil {
		senzingConfig.observers = &subject.SubjectImpl{}
	}

	// Register observer with senzingConfig and dependent services.

	err = senzingConfig.observers.RegisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 32, 1032
		return err
	}
	g2Config, g2Configmgr, err := senzingConfig.getDependentServices(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 33, 1033
		return err
	}
	err = g2Config.RegisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 34, 1034
		return err
	}
	err = g2Configmgr.RegisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 35, 1035
		return err
	}

	// Notify observers.

	go func() {
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentId, 8003, err, details)
	}()

	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (senzingConfig *SenzingConfigImpl) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil

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

		asJson, err := json.Marshal(senzingConfig)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 41, 1041
			return err
		}
		senzingConfig.log(1003, senzingConfig, string(asJson))
	}

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		traceExitMessageNumber, debugMessageNumber = 42, 1042
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}

	// Set senzingConfig log level.

	senzingConfig.logLevel = logLevelName
	err = senzingConfig.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 43, 1043
		return err
	}
	senzingConfig.isTrace = (logLevelName == logging.LevelTraceName)

	// Set log level for dependent services.

	g2Config, g2Configmgr, err := senzingConfig.getDependentServices(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 44, 1044
		return err
	}
	err = g2Config.SetLogLevel(ctx, logLevelName)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 45, 1045
		return err
	}
	err = g2Configmgr.SetLogLevel(ctx, logLevelName)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 46, 1046
		return err
	}

	// Notify observers.

	if senzingConfig.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentId, 8004, err, details)
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
func (senzingConfig *SenzingConfigImpl) SetObserverOrigin(ctx context.Context, origin string) {
	var err error = nil

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

		asJson, err := json.Marshal(senzingConfig)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 61, 1061
			return
		}
		senzingConfig.log(1004, senzingConfig, string(asJson))
	}

	// Set origin in dependent services.

	senzingConfig.observerOrigin = origin
	g2Config, g2Configmgr, err := senzingConfig.getDependentServices(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 62, 1062
		return
	}
	g2Config.SetObserverOrigin(ctx, origin)
	g2Configmgr.SetObserverOrigin(ctx, origin)

	// Notify observers.

	if senzingConfig.observers != nil {
		go func() {
			details := map[string]string{
				"origin": origin,
			}
			notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentId, 8005, err, details)
		}()
	}

}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be removed.
*/
func (senzingConfig *SenzingConfigImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil

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
			senzingConfig.traceEntry(50, observer.GetObserverId(ctx))
			defer func() {
				senzingConfig.traceExit(traceExitMessageNumber, observer.GetObserverId(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJson, err := json.Marshal(senzingConfig)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 51, 1051
			return err
		}
		senzingConfig.log(1005, senzingConfig, string(asJson))
	}

	// Unregister observers in dependencies.

	g2Config, g2Configmgr, err := senzingConfig.getDependentServices(ctx)
	err = g2Config.UnregisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 52, 1052
		return err
	}
	err = g2Configmgr.UnregisterObserver(ctx, observer)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 53, 1053
		return err
	}

	// Remove observer from this service.

	if senzingConfig.observers != nil {

		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentId, 8006, err, details)

		err = senzingConfig.observers.UnregisterObserver(ctx, observer)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 54, 1054
			return err
		}

		if !senzingConfig.observers.HasObservers(ctx) {
			senzingConfig.observers = nil
		}
	}

	return err
}
