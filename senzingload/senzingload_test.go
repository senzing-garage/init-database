package senzingload_test

import (
	"context"
	"os"
	"testing"

	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/init-database/senzingconfig"
	"github.com/senzing-garage/init-database/senzingload"
	"github.com/senzing-garage/init-database/senzingschema"
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

func TestSenzingLoadBasic_LoadURLs(test *testing.T) {
	ctx := test.Context()
	senzingLoad := getTestObject(ctx, test)
	err := senzingLoad.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
	err = senzingLoad.LoadURLs(ctx)
	require.NoError(test, err)
}

func TestSenzingConfigImpl_RegisterObserver(test *testing.T) {
	ctx := test.Context()
	senzingConfig := getTestObject(ctx, test)
	anObserver := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	err := senzingConfig.RegisterObserver(ctx, anObserver)
	require.NoError(test, err)
}

func TestSenzingConfigImpl_SetLogLevel(test *testing.T) {
	ctx := test.Context()
	senzingConfig := getTestObject(ctx, test)
	err := senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
}

func TestSenzingConfigImpl_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	senzingConfig := getTestObject(ctx, test)
	senzingConfig.SetObserverOrigin(ctx, "TestObserver")
}

func TestSenzingConfigImpl_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	senzingConfig := getTestObject(ctx, test)
	anObserver := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	err := senzingConfig.RegisterObserver(ctx, anObserver)
	require.NoError(test, err)
	err = senzingConfig.UnregisterObserver(ctx, anObserver)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, t *testing.T) *senzingload.BasicSenzingLoad {
	t.Helper()

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(t, err)

	result := &senzingload.BasicSenzingLoad{
		SenzingSettings: senzingSettings,
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

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func setup() {
	ctx := context.TODO()

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		panic(err)
	}

	senzingSchema := &senzingschema.BasicSenzingSchema{
		SenzingSettings: senzingSettings,
	}

	err = senzingSchema.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		panic(err)
	}

	err = senzingSchema.InitializeSenzing(ctx)
	if err != nil {
		panic(err)
	}

	senzingConfig := &senzingconfig.BasicSenzingConfig{
		SenzingSettings: senzingSettings,
	}

	err = senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		panic(err)
	}

	err = senzingConfig.InitializeSenzing(ctx)
	if err != nil {
		panic(err)
	}
}

func teardown() {}
