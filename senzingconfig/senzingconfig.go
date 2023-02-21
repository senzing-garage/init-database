package senzingconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/senzing/g2-sdk-go/g2api"
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
	logger                         messagelogger.MessageLoggerInterface
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

// func failOnError(msgId int, err error) {
// 	logger.Log(msgId, err)
// 	panic(err.Error())
// }

// Print error and leave program.
func errorExit(message string, err error) {
	fmt.Printf("Exit with error: %s   Error: %v\n", message, err)
	os.Exit(1)
}

func (initializer *InitializerImpl) getG2Factory(ctx context.Context) factory.SdkAbstractFactory {
	initializer.g2factorySyncOnce.Do(func() {
		if initializer.GrpcAddress != "" {
			initializer.g2factorySingleton = &factory.SdkAbstractFactoryImpl{
				GrpcAddress: initializer.GrpcAddress,
				GrpcOptions: initializer.GrpcOptions,
			}
		} else {
			initializer.g2factorySingleton = &factory.SdkAbstractFactoryImpl{}
		}
	})
	return initializer.g2factorySingleton
}

func (initializer *InitializerImpl) getG2config(ctx context.Context) g2api.G2config {
	initializer.g2configSyncOnce.Do(func() {
		var err error = nil
		initializer.g2configSingleton, err = initializer.getG2Factory(ctx).GetG2config(ctx)
		if err != nil {
			errorExit("", err)
		}
		sdkId, err := initializer.g2configSingleton.GetSdkId(ctx)
		if err != nil {
			errorExit("", err)
		}
		if sdkId == "base" {
			moduleName := initializer.SenzingModuleName
			if len(moduleName) == 0 {
				moduleName = defaultModuleName
			}
			err = initializer.g2configSingleton.Init(ctx, moduleName, initializer.SenzingEngineConfigurationJson, initializer.SenzingVerboseLogging)
			if err != nil {
				errorExit("", err)
			}
		}
	})
	return initializer.g2configSingleton
}

func (initializer *InitializerImpl) getG2configmgr(ctx context.Context) g2api.G2configmgr {
	initializer.g2configmgrSyncOnce.Do(func() {
		var err error = nil
		initializer.g2configmgrSingleton, err = initializer.getG2Factory(ctx).GetG2configmgr(ctx)
		if err != nil {
			errorExit("", err)
		}
		sdkId, err := initializer.g2configmgrSingleton.GetSdkId(ctx)
		if err != nil {
			errorExit("", err)
		}
		if sdkId == "base" {
			moduleName := initializer.SenzingModuleName
			if len(moduleName) == 0 {
				moduleName = defaultModuleName
			}
			err = initializer.g2configmgrSingleton.Init(ctx, moduleName, initializer.SenzingEngineConfigurationJson, initializer.SenzingVerboseLogging)
			if err != nil {
				errorExit("", err)
			}
		}
	})
	return initializer.g2configmgrSingleton
}

// Get the Logger singleton.
func (initializer *InitializerImpl) getLogger() messagelogger.MessageLoggerInterface {
	if initializer.logger == nil {
		initializer.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return initializer.logger
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
  - observer: The observer to be added.
*/
func (initializer *InitializerImpl) InitializeSenzingConfiguration(ctx context.Context) error {
	if initializer.isTrace {
		initializer.traceEntry(1)
	}
	entryTime := time.Now()
	var err error = nil

	// Log entry parameters.

	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, initializer.LogLevel)
	logger.Log(2000, initializer)

	// Create Senzing objects.

	g2Config := initializer.getG2config(ctx)
	g2Configmgr := initializer.getG2configmgr(ctx)

	// Determine if configuration already exists. If so, return.

	configID, err := g2Configmgr.GetDefaultConfigID(ctx)
	if err != nil {
		errorExit("", err)
	}

	if configID != 0 {
		if initializer.observers != nil {
			go func() {
				details := map[string]string{}
				initializer.notify(ctx, 8009, err, details)
			}()
		}
		logger.Log(2002, configID)
		if initializer.isTrace {
			defer initializer.traceExit(901, err, time.Since(entryTime))
		}
		return err
	}

	// Create a fresh Senzing configuration.

	configHandle, err := g2Config.Create(ctx)
	if err != nil {
		errorExit("", err)
	}

	configStr, err := g2Config.Save(ctx, configHandle)
	if err != nil {
		errorExit("", err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	configComments := fmt.Sprintf("Created by initdatabase at %s", entryTime.UTC())
	configID, err = g2Configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		errorExit("", err)
	}

	err = g2Configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		errorExit("", err)
	}

	if initializer.observers != nil {
		go func() {
			details := map[string]string{}
			initializer.notify(ctx, 8009, err, details)
		}()
	}
	if initializer.isTrace {
		defer initializer.traceExit(2, configID, err, configID, time.Since(entryTime))
	}
	return err
}

/*
The Initialize method adds the Senzing database schema and Senzing default configuration to databases.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (initializer *InitializerImpl) Initialize(ctx context.Context) error {
	if initializer.isTrace {
		initializer.traceEntry(3)
	}
	entryTime := time.Now()
	err := initializer.InitializeSenzingConfiguration(ctx)
	if err != nil {
		errorExit("", err)
	}
	if initializer.observers != nil {
		go func() {
			details := map[string]string{}
			initializer.notify(ctx, 8009, err, details)
		}()
	}
	if initializer.isTrace {
		defer initializer.traceExit(4, err, time.Since(entryTime))
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
		initializer.traceEntry(5, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	if initializer.observers == nil {
		initializer.observers = &subject.SubjectImpl{}
	}
	err := initializer.observers.RegisterObserver(ctx, observer)
	initializer.getG2config(ctx).RegisterObserver(ctx, observer)
	initializer.getG2configmgr(ctx).RegisterObserver(ctx, observer)
	if initializer.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverId(ctx),
			}
			initializer.notify(ctx, 8008, err, details)
		}()
	}
	if initializer.isTrace {
		defer initializer.traceExit(6, observer.GetObserverId(ctx), err, time.Since(entryTime))
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
		initializer.traceEntry(7, logLevel)
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
			initializer.notify(ctx, 8009, err, details)
		}()
	}
	if initializer.isTrace {
		defer initializer.traceExit(8, logLevel, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (initializer *InitializerImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	if initializer.isTrace {
		initializer.traceEntry(9, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	var err error = nil
	initializer.getG2config(ctx).UnregisterObserver(ctx, observer)
	initializer.getG2configmgr(ctx).UnregisterObserver(ctx, observer)
	if initializer.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		initializer.notify(ctx, 8013, err, details)
	}
	err = initializer.observers.UnregisterObserver(ctx, observer)
	if !initializer.observers.HasObservers(ctx) {
		initializer.observers = nil
	}
	if initializer.isTrace {
		defer initializer.traceExit(10, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}
