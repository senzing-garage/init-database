package senzingschema_test

import (
	"testing"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/init-database/senzingschema"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSenzingSchemaImpl_InitializeSenzing(test *testing.T) {
	ctx := test.Context()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)

	testObject := &senzingschema.BasicSenzingSchema{
		SenzingSettings: senzingSettings,
	}
	err = testObject.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
	err = testObject.InitializeSenzing(ctx)
	require.NoError(test, err)
}

func TestSenzingSchemaImpl_RegisterObserver(test *testing.T) {
	ctx := test.Context()
	observer1 := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)

	testObject := &senzingschema.BasicSenzingSchema{
		SenzingSettings: senzingSettings,
	}
	err = testObject.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
	err = testObject.RegisterObserver(ctx, observer1)
	require.NoError(test, err)
	err = testObject.InitializeSenzing(ctx)
	require.NoError(test, err)
}

func TestSenzingSchemaImpl_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)

	testObject := &senzingschema.BasicSenzingSchema{
		SenzingSettings: senzingSettings,
	}
	testObject.SetObserverOrigin(ctx, "TestObserver")
}

func TestSenzingSchemaImpl_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	observer1 := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)

	testObject := &senzingschema.BasicSenzingSchema{
		SenzingSettings: senzingSettings,
	}
	err = testObject.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
	err = testObject.RegisterObserver(ctx, observer1)
	require.NoError(test, err)
	err = testObject.InitializeSenzing(ctx)
	require.NoError(test, err)
	err = testObject.UnregisterObserver(ctx, observer1)
	require.NoError(test, err)
}
