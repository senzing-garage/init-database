package senzingconfig

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-observing/notifier"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
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
	observers                      subject.Subject
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var defaultModuleName string = "init-database"

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (senzingConfig *SenzingConfigImpl) getLogger() logging.LoggingInterface {
	var err error = nil
	if senzingConfig.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		senzingConfig.logger, err = logging.NewSenzingToolsLogger(ProductId, IdMessages, options...)
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

// Trace method entry.
func (senzingConfig *SenzingConfigImpl) traceEntry(messageNumber int, details ...interface{}) {
	senzingConfig.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (senzingConfig *SenzingConfigImpl) traceExit(messageNumber int, details ...interface{}) {
	senzingConfig.getLogger().Log(messageNumber, details...)
}

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
			if err != nil {
				return
			}
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
			if err != nil {
				return
			}
		}
	})
	return senzingConfig.g2configmgrSingleton, err
}

// Add datasources to Senzing configuration.
func (senzingConfig *SenzingConfigImpl) addDatasources(ctx context.Context, g2Config g2api.G2config, configHandle uintptr) error {
	var err error = nil
	for _, datasource := range senzingConfig.DataSources {
		inputJson := `{"DSRC_CODE": "` + datasource + `"}`
		_, err = g2Config.AddDataSource(ctx, configHandle, inputJson)
		if err != nil {
			return err
		}
		senzingConfig.log(2003, datasource)
	}
	return err
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Initialize method adds the Senzing default configuration to databases.

Input
  - ctx: A context to control lifecycle.
*/
func (senzingConfig *SenzingConfigImpl) Initialize(ctx context.Context) error {
	if senzingConfig.isTrace {
		senzingConfig.traceEntry(1)
	}
	entryTime := time.Now()

	// Log entry parameters.

	senzingConfig.log(1000, senzingConfig)

	// Create Senzing objects.

	g2Config, err := senzingConfig.getG2config(ctx)
	if err != nil {
		return err
	}
	g2Configmgr, err := senzingConfig.getG2configmgr(ctx)
	if err != nil {
		return err
	}

	// Determine if configuration already exists. If so, return.

	configID, err := g2Configmgr.GetDefaultConfigID(ctx)
	if err != nil {
		return err
	}
	if configID != 0 {
		if senzingConfig.observers != nil {
			go func() {
				details := map[string]string{}
				notifier.Notify(ctx, senzingConfig.observers, ProductId, 8001, err, details)
			}()
		}
		senzingConfig.log(2002, configID)
		if senzingConfig.isTrace {
			defer senzingConfig.traceExit(901, err, time.Since(entryTime))
		}
		return err
	}

	// Create a fresh Senzing configuration.

	configHandle, err := g2Config.Create(ctx)
	if err != nil {
		return err
	}

	// If requested, add DataSources to fresh Senzing configuration.

	if len(senzingConfig.DataSources) > 0 {
		err = senzingConfig.addDatasources(ctx, g2Config, configHandle)
		if err != nil {
			return err
		}
	}

	// Create a JSON string from the in-memory configuration.

	configStr, err := g2Config.Save(ctx, configHandle)
	if err != nil {
		return err
	}

	// Persist the Senzing configuration to the Senzing repository and set as default configuration.

	configComments := fmt.Sprintf("Created by init-database at %s", entryTime.Format(time.RFC3339Nano))
	configID, err = g2Configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return err
	}
	err = g2Configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return err
	}

	// Epilog.

	senzingConfig.log(2004, configID, configComments)
	if senzingConfig.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, senzingConfig.observers, ProductId, 8002, err, details)
		}()
	}
	if senzingConfig.isTrace {
		defer senzingConfig.traceExit(2, err, configID, time.Since(entryTime))
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
	if senzingConfig.isTrace {
		senzingConfig.traceEntry(3, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	if senzingConfig.observers == nil {
		senzingConfig.observers = &subject.SubjectImpl{}
	}

	// Register observer with senzingConfig.

	err := senzingConfig.observers.RegisterObserver(ctx, observer)
	if err != nil {
		return err
	}

	// Register observer with dependent stucts.

	g2Config, err := senzingConfig.getG2config(ctx)
	if err != nil {
		return err
	}
	g2Configmgr, err := senzingConfig.getG2configmgr(ctx)
	if err != nil {
		return err
	}
	for _, observer := range senzingConfig.observers.GetObservers(ctx) {
		err = g2Config.RegisterObserver(ctx, observer)
		if err != nil {
			return err
		}
		err = g2Configmgr.RegisterObserver(ctx, observer)
		if err != nil {
			return err
		}
	}

	// Epilog.

	details := map[string]string{
		"observerID": observer.GetObserverId(ctx),
	}
	notifier.Notify(ctx, senzingConfig.observers, ProductId, 8003, err, details)
	if senzingConfig.isTrace {
		defer senzingConfig.traceExit(4, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (senzingConfig *SenzingConfigImpl) SetLogLevel(ctx context.Context, logLevelName string) error {
	if senzingConfig.isTrace {
		senzingConfig.traceEntry(5, logLevelName)
	}
	entryTime := time.Now()
	var err error = nil
	if logging.IsValidLogLevelName(logLevelName) {
		senzingConfig.logLevel = logLevelName
		senzingConfig.getLogger().SetLogLevel(logLevelName)
		senzingConfig.isTrace = (logLevelName == logging.LevelTraceName)
		if senzingConfig.observers != nil {
			go func() {
				details := map[string]string{
					"logLevelName": logLevelName,
				}
				notifier.Notify(ctx, senzingConfig.observers, ProductId, 8004, err, details)
			}()
		}
	} else {
		err = fmt.Errorf("invalid error level: %s", logLevelName)
	}
	if senzingConfig.isTrace {
		defer senzingConfig.traceExit(6, logLevelName, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be removed.
*/
func (senzingConfig *SenzingConfigImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	if senzingConfig.isTrace {
		senzingConfig.traceEntry(7, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	var err error = nil
	// senzingConfig.getG2config(ctx).UnregisterObserver(ctx, observer)
	// senzingConfig.getG2configmgr(ctx).UnregisterObserver(ctx, observer)
	if senzingConfig.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, senzingConfig.observers, ProductId, 8005, err, details)
	}
	err = senzingConfig.observers.UnregisterObserver(ctx, observer)
	if !senzingConfig.observers.HasObservers(ctx) {
		senzingConfig.observers = nil
	}
	if senzingConfig.isTrace {
		defer senzingConfig.traceExit(8, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}
