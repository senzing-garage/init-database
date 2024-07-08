package senzingconfig

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/init-database/senzingschema"
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

func TestSenzingConfigImpl_InitializeSenzing_withDatasources(test *testing.T) {
	ctx := context.TODO()
	senzingConfig := getTestObject(ctx, test)
	senzingConfig.DataSources = []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	err := senzingConfig.InitializeSenzing(ctx)
	require.NoError(test, err)
}

func TestSenzingConfigImpl_InitializeSenzing(test *testing.T) {
	ctx := context.TODO()
	senzingConfig := getTestObject(ctx, test)
	err := senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
	err = senzingConfig.InitializeSenzing(ctx)
	require.NoError(test, err)
}

func TestSenzingConfigImpl_RegisterObserver(test *testing.T) {
	ctx := context.TODO()
	senzingConfig := getTestObject(ctx, test)
	anObserver := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	err := senzingConfig.RegisterObserver(ctx, anObserver)
	require.NoError(test, err)
}

func TestSenzingConfigImpl_SetLogLevel(test *testing.T) {
	ctx := context.TODO()
	senzingConfig := getTestObject(ctx, test)
	err := senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
}

func TestSenzingConfigImpl_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	senzingConfig := getTestObject(ctx, test)
	senzingConfig.SetObserverOrigin(ctx, "TestObserver")
}

func TestSenzingConfigImpl_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
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

func getTestObject(ctx context.Context, test *testing.T) *BasicSenzingConfig {
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)
	result := &BasicSenzingConfig{
		SenzingSettings: senzingSettings,
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

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	ctx := context.TODO()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &senzingschema.BasicSenzingSchema{
		SenzingSettings: senzingSettings,
	}
	err = senzingSchema.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	err = senzingSchema.InitializeSenzing(ctx)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func teardown() error {
	var err error
	return err
}
