package senzingload

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

// BasicSenzingLoad is the default implementation of the SenzingLoad interface.
type BasicSenzingLoad struct {
	JSONURLs              []string          `json:"jsonUrls,omitempty"`
	GrpcDialOptions       []grpc.DialOption `json:"grpcDialOptions,omitempty"`
	GrpcTarget            string            `json:"grpcTarget,omitempty"`
	SenzingInstanceName   string            `json:"senzingInstanceName,omitempty"`
	SenzingSettings       string            `json:"senzingSettings,omitempty"`
	SenzingVerboseLogging int64             `json:"senzingVerboseLogging,omitempty"`

	isTrace                    bool
	logger                     logging.Logging
	logLevel                   string
	observerOrigin             string
	observers                  subject.Subject
	szAbstractFactorySingleton senzing.SzAbstractFactory
	szAbstractFactorySyncOnce  sync.Once
}

type record struct {
	DataSource string `json:"DATA_SOURCE"`
	ID         string `json:"RECORD_ID"`
}

const timeoutInMinutes = 15 // Fifteen minutes is just a guess.

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
The LoadURLs method adds records from URLs to files of JSON lines.
It also process redo records before returning.

Input
  - ctx: A context to control lifecycle.
*/
func (senzingLoad *BasicSenzingLoad) LoadURLs(ctx context.Context) error {
	var err error

	entryTime := time.Now()

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 29

	if senzingLoad.getLogger().IsDebug() {
		// If DEBUG, log error exit.
		defer func() {
			if debugMessageNumber > 0 {
				senzingLoad.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingLoad.getLogger().IsTrace() {
			senzingLoad.traceEntry(10)

			defer func() { senzingLoad.traceExit(traceExitMessageNumber, err, time.Since(entryTime)) }()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingLoad)
		if err != nil {
			return wraperror.Errorf(err, "json.Marshal: %v", senzingLoad)
		}

		senzingLoad.log(1001, senzingLoad, string(asJSON))
	}

	// Create Senzing objects.

	szAbstractFactory := senzingLoad.getAbstractFactory(ctx)

	defer func() { assertNoError(szAbstractFactory.Close(ctx), "Error on szAbstractFactory.Close()") }()

	// Process each URL of JSON lines.

	err = senzingLoad.processRecords(ctx, szAbstractFactory)
	if err != nil {
		return wraperror.Errorf(err, "senzingLoad.processRecords")
	}

	// Process redo records.

	err = senzingLoad.processRedoRecords(ctx, szAbstractFactory)
	if err != nil {
		return wraperror.Errorf(err, "senzingLoad.processRedoRecords")
	}

	// Notify observers.

	if senzingLoad.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, senzingLoad.observers, senzingLoad.observerOrigin, ComponentID, 8002, nil, details)
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
func (senzingLoad *BasicSenzingLoad) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if observer == nil {
		return err
	}

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 39

	if senzingLoad.getLogger().IsDebug() {
		// If DEBUG, log error exit.
		defer func() {
			if debugMessageNumber > 0 {
				senzingLoad.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingLoad.getLogger().IsTrace() {
			entryTime := time.Now()

			senzingLoad.traceEntry(30, observer.GetObserverID(ctx))

			defer func() {
				senzingLoad.traceExit(traceExitMessageNumber, observer.GetObserverID(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingLoad)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 31, 1031

			return wraperror.Errorf(err, "json.Marshal: %v", senzingLoad)
		}

		senzingLoad.log(1002, senzingLoad, string(asJSON))
	}

	// Create empty list of observers.

	if senzingLoad.observers == nil {
		senzingLoad.observers = &subject.SimpleSubject{}
	}

	// Notify observers.

	go func() {
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, senzingLoad.observers, senzingLoad.observerOrigin, ComponentID, 8003, err, details)
	}()

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (senzingLoad *BasicSenzingLoad) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 49

	if senzingLoad.getLogger().IsDebug() {
		// If DEBUG, log error exit.
		defer func() {
			if debugMessageNumber > 0 {
				senzingLoad.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingLoad.getLogger().IsTrace() {
			entryTime := time.Now()

			senzingLoad.traceEntry(40, logLevelName)

			defer func() {
				senzingLoad.traceExit(traceExitMessageNumber, logLevelName, err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingLoad)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 41, 1041

			return wraperror.Errorf(err, "json.Marshal: %v", senzingLoad)
		}

		senzingLoad.log(1003, senzingLoad, string(asJSON))
	}

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		traceExitMessageNumber, debugMessageNumber = 42, 1042

		return wraperror.Errorf(errForPackage, "invalid error level: %s", logLevelName)
	}

	// Set senzingLoad log level.

	senzingLoad.logLevel = logLevelName

	err = senzingLoad.getLogger().SetLogLevel(logLevelName)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 43, 1043

		return wraperror.Errorf(err, "SetLogLevel: %s", logLevelName)
	}

	senzingLoad.isTrace = (logLevelName == logging.LevelTraceName)

	// Notify observers.

	if senzingLoad.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, senzingLoad.observers, senzingLoad.observerOrigin, ComponentID, 8004, err, details)
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
func (senzingLoad *BasicSenzingLoad) SetObserverOrigin(ctx context.Context, origin string) {
	var err error

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 69

	if senzingLoad.getLogger().IsDebug() {
		// If DEBUG, log error exit.
		defer func() {
			if debugMessageNumber > 0 {
				senzingLoad.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingLoad.getLogger().IsTrace() {
			entryTime := time.Now()

			senzingLoad.traceEntry(60, origin)

			defer func() {
				senzingLoad.traceExit(traceExitMessageNumber, origin, err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingLoad)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 61, 1061

			return
		}

		senzingLoad.log(1004, senzingLoad, string(asJSON))
	}

	// Notify observers.

	if senzingLoad.observers != nil {
		go func() {
			details := map[string]string{
				"origin": origin,
			}
			notifier.Notify(ctx, senzingLoad.observers, senzingLoad.observerOrigin, ComponentID, 8005, err, details)
		}()
	}
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be removed.
*/
func (senzingLoad *BasicSenzingLoad) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if observer == nil {
		return err
	}

	// Prolog.

	debugMessageNumber := 0
	traceExitMessageNumber := 59

	if senzingLoad.getLogger().IsDebug() {
		// If DEBUG, log error exit.
		defer func() {
			if debugMessageNumber > 0 {
				senzingLoad.debug(debugMessageNumber, err)
			}
		}()

		// If TRACE, Log on entry/exit.

		if senzingLoad.getLogger().IsTrace() {
			entryTime := time.Now()

			senzingLoad.traceEntry(50, observer.GetObserverID(ctx))

			defer func() {
				senzingLoad.traceExit(traceExitMessageNumber, observer.GetObserverID(ctx), err, time.Since(entryTime))
			}()
		}

		// If DEBUG, log input parameters. Must be done after establishing DEBUG and TRACE logging.

		asJSON, err := json.Marshal(senzingLoad)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 51, 1051

			return wraperror.Errorf(err, "json.Marshal: %v", senzingLoad)
		}

		senzingLoad.log(1005, senzingLoad, string(asJSON))
	}

	// Remove observer from this service.

	if senzingLoad.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, senzingLoad.observers, senzingLoad.observerOrigin, ComponentID, 8006, err, details)

		err = senzingLoad.observers.UnregisterObserver(ctx, observer)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 54, 1054

			return wraperror.Errorf(err, "UnregisterObserver")
		}

		if !senzingLoad.observers.HasObservers(ctx) {
			senzingLoad.observers = nil
		}
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (senzingLoad *BasicSenzingLoad) getLogger() logging.Logging {
	var err error

	if senzingLoad.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: OptionCallerSkip4},
			logging.OptionMessageFields{Value: []string{"id", "text"}},
		}

		senzingLoad.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}

	return senzingLoad.logger
}

// Log message.
func (senzingLoad *BasicSenzingLoad) log(messageNumber int, details ...interface{}) {
	senzingLoad.getLogger().Log(messageNumber, details...)
}

// Debug.
func (senzingLoad *BasicSenzingLoad) debug(messageNumber int, details ...interface{}) {
	details = append(details, debugOptions...)
	senzingLoad.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (senzingLoad *BasicSenzingLoad) traceEntry(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	senzingLoad.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (senzingLoad *BasicSenzingLoad) traceExit(messageNumber int, details ...interface{}) {
	details = append(details, traceOptions...)
	senzingLoad.getLogger().Log(messageNumber, details...)
}

// --- Dependent services -----------------------------------------------------

// Create an abstract factory singleton and return it.
func (senzingLoad *BasicSenzingLoad) getAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	var err error

	_ = ctx // ctx not used, yet.

	senzingLoad.szAbstractFactorySyncOnce.Do(func() {
		if len(senzingLoad.GrpcTarget) == 0 {
			senzingLoad.szAbstractFactorySingleton, err = senzingLoad.buildSzAbstractFactory()
			if err != nil {
				panic(err)
			}
		} else {
			grpcConnection, err := grpc.NewClient(senzingLoad.GrpcTarget, senzingLoad.GrpcDialOptions...)
			if err != nil {
				panic(err)
			}

			senzingLoad.szAbstractFactorySingleton, err = szfactorycreator.CreateGrpcAbstractFactory(grpcConnection)
			if err != nil {
				panic(err)
			}
		}
	})

	return senzingLoad.szAbstractFactorySingleton
}

func (senzingLoad *BasicSenzingLoad) processRecords(
	ctx context.Context,
	szAbstractFactory senzing.SzAbstractFactory,
) error {
	var (
		jsonRecord record
		err        error
	)

	// Use timeout in ctx.

	// Get an szEngine.

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		return wraperror.Errorf(err, "CreateEngine")
	}

	defer func() { assertNoError(szEngine.Destroy(ctx), "Error on szEngine.Destroy()") }()

	httpClient := &http.Client{
		Timeout: timeoutInMinutes * time.Minute,
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, timeoutInMinutes*time.Minute)
	defer cancel() // Ensure the context is canceled when main exits

	for _, jsonURL := range senzingLoad.JSONURLs {
		select {
		case <-ctxTimeout.Done():
			return wraperror.Errorf(ctx.Err(), "HTTP Timeout")
		default:
			senzingLoad.log(3001, jsonURL)

			// Download file from URL.

			httpRequest, errNewRequestWithContext := http.NewRequestWithContext(
				ctxTimeout,
				http.MethodGet,
				jsonURL,
				nil,
			)
			if errNewRequestWithContext != nil {
				return wraperror.Errorf(errNewRequestWithContext, "http.NewRequestWithContext")
			}

			httpResponse, errDo := httpClient.Do(httpRequest)
			if errDo != nil {
				return wraperror.Errorf(errDo, "httpClient.Do")
			}

			if httpResponse.StatusCode != http.StatusOK {
				errBodyClose := httpResponse.Body.Close()

				return wraperror.Errorf(
					errForPackage,
					fmt.Sprintf(
						"Received non-OK HTTP status: %d; Error for httpResponse.Body.Close: %v",
						httpResponse.StatusCode,
						errBodyClose,
					),
				)
			}

			// Process HTTP response body in-memory. That is, do not write to file.

			jsonLineCount := 0

			scanner := bufio.NewScanner(httpResponse.Body)
			for scanner.Scan() {
				jsonLineCount++
				line := scanner.Bytes()

				// Tricky code: The following code assumes that the JSON contains "DATA_SOURCE" and "RECORD_ID" JSON keys.

				err = json.Unmarshal(line, &jsonRecord)
				if err != nil {
					errBodyClose := httpResponse.Body.Close()

					return wraperror.Errorf(
						err,
						fmt.Sprintf("Scanning:  %s; Error for httpResponse.Body.Close: %v", string(line), errBodyClose),
					)
				}

				_, err = szEngine.AddRecord(
					ctxTimeout,
					jsonRecord.DataSource,
					jsonRecord.ID,
					string(line),
					senzing.SzNoFlags,
				)
				if err != nil {
					errBodyClose := httpResponse.Body.Close()

					return wraperror.Errorf(
						err,
						fmt.Sprintf(
							"szEngine.AddRecord DataSource: %s; RecordID: %s; Error for httpResponse.Body.Close: %v",
							jsonRecord.DataSource,
							jsonRecord.ID,
							errBodyClose,
						),
					)
				}
			}

			if err := scanner.Err(); err != nil {
				return wraperror.Errorf(err, "Scanning")
			}

			err = httpResponse.Body.Close()
			if err != nil {
				return wraperror.Errorf(err, "Error for httpResponse.Body.Close")
			}

			senzingLoad.log(2002, jsonLineCount, jsonURL)
		}
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

func (senzingLoad *BasicSenzingLoad) processRedoRecords(
	ctx context.Context,
	szAbstractFactory senzing.SzAbstractFactory,
) error {
	var err error

	// Get an szEngine.

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		return wraperror.Errorf(err, "CreateEngine")
	}

	defer func() { assertNoError(szEngine.Destroy(ctx), "Error on szEngine.Destroy()") }()

	// Process Senzing Redo Records.

	redoRecordCount := 0

	for {
		redoRecord, err := szEngine.GetRedoRecord(ctx)
		if err != nil {
			return wraperror.Errorf(err, "szEngine.GetRedoRecord()")
		}

		if len(redoRecord) == 0 {
			break
		}

		redoRecordCount++

		_, err = szEngine.ProcessRedoRecord(ctx, redoRecord, senzing.SzNoFlags)
		if err != nil {
			return wraperror.Errorf(err, "szEngine.ProcessRedoRecord()")
		}
	}

	senzingLoad.log(2003, redoRecordCount)

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// --- Misc -------------------------------------------------------------------

func (senzingLoad *BasicSenzingLoad) buildSzAbstractFactory() (senzing.SzAbstractFactory, error) {
	var (
		err    error
		result senzing.SzAbstractFactory
	)

	senzingSettings := senzingLoad.SenzingSettings

	senzingInstanceName := senzingLoad.SenzingInstanceName
	if len(senzingInstanceName) == 0 {
		senzingInstanceName = fmt.Sprintf("senzing init-database at %s", time.Now())
	}

	result, err = szfactorycreator.CreateCoreAbstractFactory(
		senzingInstanceName,
		senzingSettings,
		senzingLoad.SenzingVerboseLogging,
		senzing.SzInitializeWithDefaultConfiguration,
	)

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func assertNoError(err error, message string) {
	if err != nil {
		log.Fatalf("Error: %s; err = %v", message, err)
	}
}
