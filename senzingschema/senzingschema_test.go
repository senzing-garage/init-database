package senzingschema

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
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
	var err error = nil
	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSenzingSchemaImpl_InitializeSenzing(test *testing.T) {
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := settings.BuildSimpleSettingsUsingEnvVars()
	testError(test, err)
	testObject := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	testObject.SetLogLevel(ctx, logging.LevelInfoName)
	testObject.InitializeSenzing(ctx)
}

func TestSenzingSchemaImpl_RegisterObserver(test *testing.T) {
	ctx := context.TODO()
	observer1 := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := settings.BuildSimpleSettingsUsingEnvVars()
	testError(test, err)
	testObject := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	testObject.SetLogLevel(ctx, logging.LevelInfoName)
	testObject.RegisterObserver(ctx, observer1)
	testObject.InitializeSenzing(ctx)
}

func TestSenzingSchemaImpl_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := settings.BuildSimpleSettingsUsingEnvVars()
	testError(test, err)
	testObject := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	testObject.SetObserverOrigin(ctx, "TestObserver")
}

func TestSenzingSchemaImpl_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	observer1 := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := settings.BuildSimpleSettingsUsingEnvVars()
	testError(test, err)
	testObject := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	testObject.SetLogLevel(ctx, logging.LevelInfoName)
	testObject.RegisterObserver(ctx, observer1)
	testObject.InitializeSenzing(ctx)
	testObject.UnregisterObserver(ctx, observer1)
}
