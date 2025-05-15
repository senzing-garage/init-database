//go:build linux

package senzingconfig_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/init-database/senzingconfig"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleBasicSenzingConfig_InitializeSenzing_withDatasources() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}

	senzingConfig := &senzingconfig.BasicSenzingConfig{
		SenzingSettings: senzingSettings,
		DataSources:     []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"},
	}

	err = senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}

	_ = senzingConfig.InitializeSenzing(ctx)
	// Output:
}

func ExampleBasicSenzingConfig_InitializeSenzing() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}

	senzingConfig := &senzingconfig.BasicSenzingConfig{
		SenzingSettings: senzingSettings,
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

func ExampleBasicSenzingConfig_RegisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}

	senzingConfig := &senzingconfig.BasicSenzingConfig{
		SenzingSettings: senzingSettings,
	}

	err = senzingConfig.RegisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleBasicSenzingConfig_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}

	senzingConfig := &senzingconfig.BasicSenzingConfig{
		SenzingSettings: senzingSettings,
	}

	err = senzingConfig.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleBasicSenzingConfig_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}

	senzingConfig := &senzingconfig.BasicSenzingConfig{
		SenzingSettings: senzingSettings,
	}
	senzingConfig.SetObserverOrigin(ctx, "TestObserver")
	// Output:
}

func ExampleBasicSenzingConfig_UnregisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingconfig/senzingconfig_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}

	senzingConfig := &senzingconfig.BasicSenzingConfig{
		SenzingSettings: senzingSettings,
	}

	err = senzingConfig.RegisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Println(err)
	}

	err = senzingConfig.UnregisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
