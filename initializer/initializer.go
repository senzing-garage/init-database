package initializer

import (
	"context"
	"time"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/notifier"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
	"github.com/senzing/initdatabase/senzingconfig"
	"github.com/senzing/initdatabase/senzingschema"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// InitializerImpl is the default implementation of the Initializer interface.
type InitializerImpl struct {
	GrpcAddress                    string
	GrpcOptions                    []grpc.DialOption
	isTrace                        bool
	logLevel                       logger.Level
	messageLogger                  messagelogger.MessageLoggerInterface
	observers                      subject.Subject
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (initializerImpl *InitializerImpl) getLogger() messagelogger.MessageLoggerInterface {
	if initializerImpl.messageLogger == nil {
		initializerImpl.messageLogger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, initializerImpl.logLevel)
	}
	return initializerImpl.messageLogger
}

// Trace method entry.
func (initializerImpl *InitializerImpl) traceEntry(errorNumber int, details ...interface{}) {
	initializerImpl.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (initializerImpl *InitializerImpl) traceExit(errorNumber int, details ...interface{}) {
	initializerImpl.getLogger().Log(errorNumber, details...)
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
	if initializerImpl.isTrace {
		initializerImpl.traceEntry(1)
	}
	entryTime := time.Now()

	// Log entry parameters.

	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, initializerImpl.logLevel)
	logger.Log(2000, initializerImpl)

	// Create senzingSchema for initializing Senzing schema.

	senzingSchema := &senzingschema.SenzingSchemaImpl{
		SenzingEngineConfigurationJson: initializerImpl.SenzingEngineConfigurationJson,
	}
	senzingSchema.SetLogLevel(ctx, initializerImpl.logLevel)

	// Create senzingConfig for initializing Senzing configuration.

	senzingConfig := &senzingconfig.SenzingConfigImpl{
		GrpcAddress:                    initializerImpl.GrpcAddress,
		GrpcOptions:                    initializerImpl.GrpcOptions,
		SenzingEngineConfigurationJson: initializerImpl.SenzingEngineConfigurationJson,
		SenzingModuleName:              initializerImpl.SenzingModuleName,
		SenzingVerboseLogging:          initializerImpl.SenzingVerboseLogging,
	}
	senzingConfig.SetLogLevel(ctx, initializerImpl.logLevel)

	// Add observers to structs.

	if initializerImpl.observers != nil {
		for _, observer := range initializerImpl.observers.GetObservers(ctx) {
			err = senzingConfig.RegisterObserver(ctx, observer)
			if err != nil {
				return err
			}
			err = senzingSchema.RegisterObserver(ctx, observer)
			if err != nil {
				return err
			}
		}
	}

	// Perform initialization.

	err = senzingSchema.Initialize(ctx)
	if err != nil {
		return err
	}
	err = senzingConfig.Initialize(ctx)
	if err != nil {
		return err
	}

	// Epilog.

	if initializerImpl.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, initializerImpl.observers, ProductId, 8001, err, details)
		}()
	}
	if initializerImpl.isTrace {
		defer initializerImpl.traceExit(2, err, time.Since(entryTime))
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
	if initializerImpl.isTrace {
		initializerImpl.traceEntry(3, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	if initializerImpl.observers == nil {
		initializerImpl.observers = &subject.SubjectImpl{}
	}
	err := initializerImpl.observers.RegisterObserver(ctx, observer)
	details := map[string]string{
		"observerID": observer.GetObserverId(ctx),
	}
	notifier.Notify(ctx, initializerImpl.observers, ProductId, 8002, err, details)
	if initializerImpl.isTrace {
		defer initializerImpl.traceExit(4, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (initializerImpl *InitializerImpl) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if initializerImpl.isTrace {
		initializerImpl.traceEntry(5, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	initializerImpl.logLevel = logLevel
	initializerImpl.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	initializerImpl.isTrace = (initializerImpl.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if initializerImpl.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logger.LevelToTextMap[logLevel],
			}
			notifier.Notify(ctx, initializerImpl.observers, ProductId, 8003, err, details)
		}()
	}
	if initializerImpl.isTrace {
		defer initializerImpl.traceExit(6, logLevel, err, time.Since(entryTime))
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
	if initializerImpl.isTrace {
		initializerImpl.traceEntry(7, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	var err error = nil
	if initializerImpl.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, initializerImpl.observers, ProductId, 8004, err, details)
	}
	err = initializerImpl.observers.UnregisterObserver(ctx, observer)
	if !initializerImpl.observers.HasObservers(ctx) {
		initializerImpl.observers = nil
	}
	if initializerImpl.isTrace {
		defer initializerImpl.traceExit(8, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}
