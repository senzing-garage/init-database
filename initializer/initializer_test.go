package initializer

import (
	"context"
	"testing"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestInitializerImpl_Initialize(test *testing.T) {
	ctx := context.TODO()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)
	testObject := &BasicInitializer{
		SenzingSettings: senzingSettings,
	}
	err = testObject.SetLogLevel(ctx, logging.LevelInfoName)
	require.NoError(test, err)
	err = testObject.Initialize(ctx)
	require.NoError(test, err)
}

func TestInitializerImpl_RegisterObserver(test *testing.T) {
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

func TestInitializerImpl_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(test, err)
	testObject := &BasicInitializer{
		SenzingSettings: senzingSettings,
	}
	testObject.SetObserverOrigin(ctx, "TestObserver")
}

// func TestInitializerImpl_UnregisterObserver(test *testing.T) {
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
