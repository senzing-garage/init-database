package senzingconfig

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
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
	logger                         messagelogger.MessageLoggerInterface
	logLevel                       logger.Level
	observers                      subject.Subject
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var defaultModuleName string = "initdatabase"

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (senzingConfig *SenzingConfigImpl) getLogger() messagelogger.MessageLoggerInterface {
	if senzingConfig.logger == nil {
		senzingConfig.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, senzingConfig.logLevel)
	}
	return senzingConfig.logger
}

// Trace method entry.
func (senzingConfig *SenzingConfigImpl) traceEntry(errorNumber int, details ...interface{}) {
	senzingConfig.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (senzingConfig *SenzingConfigImpl) traceExit(errorNumber int, details ...interface{}) {
	senzingConfig.getLogger().Log(errorNumber, details...)
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
		senzingConfig.getLogger().Log(2003, datasource)
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

	senzingConfig.getLogger().Log(1000, senzingConfig)

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
		senzingConfig.getLogger().Log(2002, configID)
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

	configComments := fmt.Sprintf("Created by initdatabase at %s", entryTime.UTC())
	configID, err = g2Configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return err
	}
	err = g2Configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return err
	}

	// Epilog.

	senzingConfig.getLogger().Log(2004, configID, configComments)
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
func (senzingConfig *SenzingConfigImpl) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if senzingConfig.isTrace {
		senzingConfig.traceEntry(5, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	senzingConfig.logLevel = logLevel
	senzingConfig.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	senzingConfig.isTrace = (senzingConfig.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if senzingConfig.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logger.LevelToTextMap[logLevel],
			}
			notifier.Notify(ctx, senzingConfig.observers, ProductId, 8004, err, details)
		}()
	}
	if senzingConfig.isTrace {
		defer senzingConfig.traceExit(6, logLevel, err, time.Since(entryTime))
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
