package initializer

import (
	"context"
	"testing"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/stretchr/testify/require"
)

const (
	observerOrigin = "init-database observer"
)

var (
	logLevel          = helper.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestBasicInitializer_Initialize(test *testing.T) {
	ctx := context.TODO()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
	err = testObject.Initialize(ctx)
	require.NoError(test, err)
}

func TestBasicInitializer_RegisterObserver(test *testing.T) {
	ctx := context.TODO()
	observer1 := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)
	testObject := &BasicInitializer{
		SenzingSettings: senzingSettings,
	}
	err = testObject.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
	err = testObject.RegisterObserver(ctx, observer1)
	require.NoError(test, err)
	err = testObject.Initialize(ctx)
	require.NoError(test, err)
}

func TestBasicInitializer_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)
	testObject := &BasicInitializer{
		SenzingSettings: senzingSettings,
	}
	testObject.SetObserverOrigin(ctx, "TestObserver")
}

// func TestBasicInitializer_UnregisterObserver(test *testing.T) {
// 	ctx := context.TODO()
// 	observer1 := &observer.NullObserver{
// 		Id:       "Observer 1",
// 		IsSilent: true,
// 	}
// 	senzingEngineConfigurationJson, err := settings.BuildSimpleSettingsUsingEnvVars()
// 	testError(test, err)
// 	testObject := &InitializerImpl{
// 		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
// 	}
// 	testObject.SetLogLevel(ctx, logging.LevelInfoName)
// 	testObject.RegisterObserver(ctx, observer1)
// 	testObject.Initialize(ctx)
// 	testObject.UnregisterObserver(ctx, observer1)
// }

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) *BasicInitializer {
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)
	result := &BasicInitializer{
		SenzingSettings: senzingSettings,
		SenzingLogLevel: logLevel,
	}

	if logLevel == "TRACE" {
		result.SetObserverOrigin(ctx, observerOrigin)
		require.NoError(test, err)
		err = result.SetLogLevel(ctx, logLevel)
		require.NoError(test, err)
		err = result.RegisterObserver(ctx, observerSingleton)
		require.NoError(test, err)
		err = result.SetLogLevel(ctx, logLevel)
		require.NoError(test, err)
	}
	return result
}
