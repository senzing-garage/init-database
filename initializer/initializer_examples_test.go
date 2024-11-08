//go:build linux

package initializer

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

func ExampleBasicInitializer_Initialize() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/initializer/initializer_examples_test.go
	ctx := context.TODO()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf(">>>>> senzingSettings: %s\n", senzingSettings)
	anInitializer := &BasicInitializer{
		SenzingSettings: senzingSettings,
	}
	err = anInitializer.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	err = anInitializer.Initialize(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleBasicInitializer_RegisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/initializer/initializer_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	anInitializer := &BasicInitializer{
		SenzingSettings: senzingSettings,
	}
	err = anInitializer.RegisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleBasicInitializer_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/initializer/initializer_examples_test.go
	ctx := context.TODO()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	anInitializer := &BasicInitializer{
		SenzingSettings: senzingSettings,
	}
	err = anInitializer.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleBasicInitializer_SetObserverOrigin() {
	ctx := context.TODO()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	anInitializer := &BasicInitializer{
		SenzingSettings: senzingSettings,
	}
	anInitializer.SetObserverOrigin(ctx, "TestObserver")
	// Output:
}

func ExampleBasicInitializer_UnregisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/initializer/initializer_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	anInitializer := &BasicInitializer{
		SenzingSettings: senzingSettings,
	}
	err = anInitializer.RegisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Println(err)
	}
	err = anInitializer.UnregisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
