package initializer_test

import (
	"context"
	"testing"

	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/init-database/initializer"
	"github.com/stretchr/testify/require"
)

const (
	observerOrigin = "init-database observer"
)

var (
	logLevel          = env.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestBasicInitializer_Initialize(test *testing.T) {
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
	err = testObject.Initialize(ctx)
	require.NoError(test, err)
}

func TestBasicInitializer_RegisterObserver(test *testing.T) {
	ctx := test.Context()
	observer1 := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)

	testObject := &initializer.BasicInitializer{
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
	ctx := test.Context()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)

	testObject := &initializer.BasicInitializer{
		SenzingSettings: senzingSettings,
	}
	testObject.SetObserverOrigin(ctx, "TestObserver")
}

func TestBasicInitializer_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	observer1 := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)

	testObject := &initializer.BasicInitializer{
		SenzingSettings: senzingSettings,
	}
	err = testObject.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
	err = testObject.RegisterObserver(ctx, observer1)
	require.NoError(test, err)
	err = testObject.Initialize(ctx)
	require.NoError(test, err)
	err = testObject.UnregisterObserver(ctx, observer1)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, t *testing.T) *initializer.BasicInitializer {
	t.Helper()

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(t, err)

	result := &initializer.BasicInitializer{
		SenzingSettings: senzingSettings,
		SenzingLogLevel: logLevel,
	}

	if logLevel == "TRACE" {
		result.SetObserverOrigin(ctx, observerOrigin)
		require.NoError(t, err)
		err = result.SetLogLevel(ctx, logLevel)
		require.NoError(t, err)
		err = result.RegisterObserver(ctx, observerSingleton)
		require.NoError(t, err)
		err = result.SetLogLevel(ctx, logLevel)
		require.NoError(t, err)
	}

	return result
}
