package senzingconfig

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/init-database/senzingschema"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func testError(test *testing.T, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
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
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &senzingschema.SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
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
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSenzingConfigImpl_InitializeSenzing_withDatasources(test *testing.T) {
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	testError(test, err)
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		DataSources:                    []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"},
	}
	err = senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	testError(test, err)
	err = senzingConfig.InitializeSenzing(ctx)
	testError(test, err)
}

func TestSenzingConfigImpl_InitializeSenzing(test *testing.T) {
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	testError(test, err)
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	testError(test, err)
	err = senzingConfig.InitializeSenzing(ctx)
	testError(test, err)
}

func TestSenzingConfigImpl_RegisterObserver(test *testing.T) {
	ctx := context.TODO()
	anObserver := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	testError(test, err)
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingConfig.RegisterObserver(ctx, anObserver)
	testError(test, err)
}

func TestSenzingConfigImpl_SetLogLevel(test *testing.T) {
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	testError(test, err)
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	testError(test, err)
}

func TestSenzingConfigImpl_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	testError(test, err)
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	senzingConfig.SetObserverOrigin(ctx, "TestObserver")
}

func TestSenzingConfigImpl_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	anObserver := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	testError(test, err)
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingConfig.RegisterObserver(ctx, anObserver)
	testError(test, err)
	err = senzingConfig.UnregisterObserver(ctx, anObserver)
	testError(test, err)
}
