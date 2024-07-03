//go:build linux

package senzingconfig

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSenzingConfigImpl_InitializeSenzing_withDatasources() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		DataSources:                    []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"},
	}
	err = senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	err = senzingConfig.InitializeSenzing(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSenzingConfigImpl_InitializeSenzing() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	err = senzingConfig.InitializeSenzing(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSenzingConfigImpl_RegisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingConfig.RegisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Print(err)
	}
	// Output:
}

func ExampleSenzingConfigImpl_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSenzingConfigImpl_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	senzingConfig.SetObserverOrigin(ctx, "TestObserver")
	// Output:
}

func ExampleSenzingConfigImpl_UnregisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingConfig.RegisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Print(err)
	}
	err = senzingConfig.UnregisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Print(err)
	}
	// Output:
}
