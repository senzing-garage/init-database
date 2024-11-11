package senzingconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-helpers/settingsparser"
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
	SenzingSettingsFile   string            `json:"senzingSettingsFile,omitempty"`
	SenzingVerboseLogging int64             `json:"senzingVerboseLogging,omitempty"`

	isTrace                    bool
	logger                     logging.Logging
	logLevel                   string
	observerOrigin             string
	observers                  subject.Subject
	szAbstractFactorySingleton senzing.SzAbstractFactory
	szAbstractFactorySyncOnce  sync.Once
	szConfigManagerSingleton   senzing.SzConfigManager
	szConfigManagerSyncOnce    sync.Once
	szConfigSingleton          senzing.SzConfig
	szConfigSyncOnce           sync.Once
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

var defaultModuleName = "init-database"

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (senzingConfig *BasicSenzingConfig) getLogger() logging.Logging {
	var err error
	if senzingConfig.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
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
	_ = ctx
	senzingConfig.szAbstractFactorySyncOnce.Do(func() {

		if len(senzingConfig.GrpcTarget) == 0 {

			senzingSettings := senzingConfig.SenzingSettings

			// Handle case of SQLite in-memory database.
			// TODO:  Refactor to different, reusable, location.

			settingsParser := settingsparser.BasicSettingsParser{
				Settings: senzingConfig.SenzingSettings,
			}

			databaseURLs, err := settingsParser.GetDatabaseURLs(ctx)
			assertNoError(err)
			if len(databaseURLs) > 1 {
				panic("Too many databaseURLs")
			}
			databaseURL := databaseURLs[0]
			parsedURL, err := url.Parse(databaseURL)
			assertNoError(err)

			if parsedURL.Scheme == "sqlite3" {
				if len(parsedURL.RawQuery) > 0 {

					configPath, err := settingsParser.GetConfigPath(ctx)
					assertNoError(err)
					resourcePath, err := settingsParser.GetResourcePath(ctx)
					assertNoError(err)
					supportPath, err := settingsParser.GetSupportPath(ctx)
					assertNoError(err)
					licenseStringBase64, err := settingsParser.GetLicenseStringBase64(ctx)
					assertNoError(err)

					szConfiguration := settings.SzConfiguration{
						Pipeline: settings.SzConfigurationPipeline{
							ConfigPath:          configPath,
							LicenseStringBase64: licenseStringBase64,
							ResourcePath:        resourcePath,
							SupportPath:         supportPath,
						},
						SQL: settings.SzConfigurationSQL{
							Connection: databaseURL,
						},
					}

					buffer := &bytes.Buffer{}
					jsonEncoder := json.NewEncoder(buffer)
					jsonEncoder.SetEscapeHTML(false)
					err = jsonEncoder.Encode(szConfiguration)
					assertNoError(err)

					senzingSettings = buffer.String()

					// senzingSettingsBytes, err := json.Marshal(szConfiguration, encOpts())
					// assertNoError(err)

					// senzingSettings := fmt.Sprintf("%s", senzingSettingsBytes)

				}
			}
			// fmt.Printf(">>>>> senzingconfig_basic.go getAbstractFactory senzingSettings: %s\n", senzingSettings)

			senzingConfig.szAbstractFactorySingleton, err = szfactorycreator.CreateCoreAbstractFactory(senzingConfig.SenzingInstanceName, senzingSettings, senzingConfig.SenzingVerboseLogging, senzing.SzInitializeWithDefaultConfiguration)
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

// Create a SzConfig singleton and return it.
func (senzingConfig *BasicSenzingConfig) getSzConfig(ctx context.Context) (senzing.SzConfig, error) {
	var err error
	senzingConfig.szConfigSyncOnce.Do(func() {
		senzingConfig.szConfigSingleton, err = senzingConfig.getAbstractFactory(ctx).CreateSzConfig(ctx)
	})
	return senzingConfig.szConfigSingleton, err
}

// Create a SzConfigManager singleton and return it.
func (senzingConfig *BasicSenzingConfig) getSzConfigmgr(ctx context.Context) (senzing.SzConfigManager, error) {
	var err error
	senzingConfig.szConfigManagerSyncOnce.Do(func() {
		senzingConfig.szConfigManagerSingleton, err = senzingConfig.getAbstractFactory(ctx).CreateSzConfigManager(ctx)
	})
	return senzingConfig.szConfigManagerSingleton, err
}

// Get dependent services: SzConfig, SzConfigManager
func (senzingConfig *BasicSenzingConfig) getDependentServices(ctx context.Context) (senzing.SzConfig, senzing.SzConfigManager, error) {
	szConfig, err := senzingConfig.getSzConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	szConfigManager, err := senzingConfig.getSzConfigmgr(ctx)
	if err != nil {
		return szConfig, nil, err
	}
	return szConfig, szConfigManager, err
}

// --- Misc -------------------------------------------------------------------

// Add datasources to Senzing configuration.
func (senzingConfig *BasicSenzingConfig) addDatasources(ctx context.Context, szConfig senzing.SzConfig, configHandle uintptr) error {
	var err error
	for _, datasource := range senzingConfig.DataSources {
		_, err = szConfig.AddDataSource(ctx, configHandle, datasource)
		if err != nil {
			return err
		}
		senzingConfig.log(2001, datasource)
	}
	return err
}

func (senzingConfig *BasicSenzingConfig) copyFile(sourceFilename string, targetFilename string) error {
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

func (senzingConfig *BasicSenzingConfig) filesAreEqual(sourceFilename string, targetFilename string) bool {
	var (
		chunkSize        = 64000
		shortCircuitExit bool
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
			switch {
			case sourceError == io.EOF && targetError == io.EOF:
				return true
			case sourceError == io.EOF || targetError == io.EOF:
				return false
			default:
				senzingConfig.log(4001, sourceFilename, targetFilename, sourceError, targetError)
			}
		}
		if !bytes.Equal(sourceBytes, targetBytes) {
			return false
		}
	}
}

func assertNoError(err error) {
	if err != nil {
		panic(err)
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
			return err
		}
		senzingConfig.log(1001, senzingConfig, string(asJSON))
	}

	// Create Senzing objects.

	szConfig, szConfigManager, err := senzingConfig.getDependentServices(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 12, 1012
		return err
	}

	// Determine if configuration already exists. If so, return.

	configID, err = szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 13, 1013
		return err
	}
	if configID != 0 {
		if senzingConfig.observers != nil {
			go func() {
				details := map[string]string{}
				notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentID, 8001, err, details)
			}()
		}
		senzingConfig.log(2002, configID)
		traceExitMessageNumber, debugMessageNumber = 14, 0 // debugMessageNumber=0 because it's not an error.
		return err
	}

	// If engine configuration file specified, swap it in.

	if len(senzingConfig.SenzingSettingsFile) > 0 {
		parsedJSON, err := settingsparser.New(senzingConfig.SenzingSettings)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 20, 1020
			return err
		}
		resourcePath, err := parsedJSON.GetResourcePath(ctx)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 21, 1021
			return err
		}

		// Compare file names.

		sourceFilename := senzingConfig.SenzingSettingsFile
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

				// Copy source file to target to "fake out" Senzing's SzEngine.Create().

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

	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 15, 1015
		return err
	}

	// If requested, add DataSources to fresh Senzing configuration.

	if len(senzingConfig.DataSources) > 0 {
		err = senzingConfig.addDatasources(ctx, szConfig, configHandle)
		if err != nil {
			traceExitMessageNumber, debugMessageNumber = 16, 1016
			return err
		}
	}

	// Create a JSON string from the in-memory configuration.

	configStr, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 17, 1017
		return err
	}

	// Persist the Senzing configuration to the Senzing repository and set as default configuration.

	configComments := fmt.Sprintf("Created by %s at %s", defaultModuleName, entryTime.Format(time.RFC3339Nano))
	configID, err = szConfigManager.AddConfig(ctx, configStr, configComments)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 18, 1018
		return err
	}
	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		traceExitMessageNumber, debugMessageNumber = 19, 1019
		return err
	}

	// Notify observers.

	senzingConfig.log(2003, configID, configComments)
	if senzingConfig.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentID, 8002, err, details)
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
			return err
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

	return err
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
			return err
		}
		senzingConfig.log(1003, senzingConfig, string(asJSON))
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

	// Notify observers.

	if senzingConfig.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, senzingConfig.observers, senzingConfig.observerOrigin, ComponentID, 8004, err, details)
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
			return err
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
			return err
		}

		if !senzingConfig.observers.HasObservers(ctx) {
			senzingConfig.observers = nil
		}
	}

	return err
}
