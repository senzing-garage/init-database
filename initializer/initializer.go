package initializer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
	"github.com/senzing/initdatabase/senzingconfig"
	"github.com/senzing/initdatabase/senzingschema"
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
	messageLogger                  messagelogger.MessageLoggerInterface
	LogLevel                       logger.Level
	observers                      subject.Subject
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (initializer *InitializerImpl) getLogger() messagelogger.MessageLoggerInterface {
	if initializer.messageLogger == nil {
		initializer.messageLogger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return initializer.messageLogger
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
Essentially it calls InitializeSenzingConfiguration() and InitializeSenzingSchema().

Input
  - ctx: A context to control lifecycle.
*/
func (initializer *InitializerImpl) Initialize(ctx context.Context) error {
	var err error = nil
	if initializer.isTrace {
		initializer.traceEntry(5)
	}
	entryTime := time.Now()

	// Initialize Senzing schema.

	senzingSchema := &senzingschema.SenzingSchemaImpl{
		LogLevel:                       logger.LevelInfo,
		SenzingEngineConfigurationJson: initializer.SenzingEngineConfigurationJson,
	}

	// FIXME:
	// for _, observer := range initializer.observers.GetObservers(ctx) {
	// 	senzingSchema.RegisterObserver(observer)
	// }
	senzingSchema.Initialize(ctx)

	// Initialize Senzing configuration.

	senzingConfig := &senzingconfig.SenzingConfigImpl{
		GrpcAddress:                    initializer.GrpcAddress,
		GrpcOptions:                    initializer.GrpcOptions,
		LogLevel:                       logger.LevelInfo,
		SenzingEngineConfigurationJson: initializer.SenzingEngineConfigurationJson,
		SenzingModuleName:              initializer.SenzingModuleName,
		SenzingVerboseLogging:          initializer.SenzingVerboseLogging,
	}

	// FIXME:
	// for _, observer := range initializer.observers.GetObservers(ctx) {
	// 	senzingSchema.RegisterObserver(observer)
	// }
	senzingConfig.Initialize(ctx)

	// Epilog.

	if initializer.observers != nil {
		go func() {
			details := map[string]string{}
			initializer.notify(ctx, 8004, err, details)
		}()
	}
	if initializer.isTrace {
		defer initializer.traceExit(6, err, time.Since(entryTime))
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
		initializer.traceEntry(7, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	if initializer.observers == nil {
		initializer.observers = &subject.SubjectImpl{}
	}
	err := initializer.observers.RegisterObserver(ctx, observer)
	if initializer.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverId(ctx),
			}
			initializer.notify(ctx, 8005, err, details)
		}()
	}
	if initializer.isTrace {
		defer initializer.traceExit(8, observer.GetObserverId(ctx), err, time.Since(entryTime))
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
		initializer.traceEntry(9, logLevel)
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
			initializer.notify(ctx, 8006, err, details)
		}()
	}
	if initializer.isTrace {
		defer initializer.traceExit(10, logLevel, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be removed.
*/
func (initializer *InitializerImpl) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	if initializer.isTrace {
		initializer.traceEntry(11, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	var err error = nil
	if initializer.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		initializer.notify(ctx, 8007, err, details)
	}
	err = initializer.observers.UnregisterObserver(ctx, observer)
	if !initializer.observers.HasObservers(ctx) {
		initializer.observers = nil
	}
	if initializer.isTrace {
		defer initializer.traceExit(12, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}
