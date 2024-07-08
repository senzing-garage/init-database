//go:build linux

package senzingschema

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

func ExampleBasicSenzingSchema_InitializeSenzing() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingschema/senzingschema_examples_test.go
	ctx := context.TODO()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &BasicSenzingSchema{
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
	// Output:
}

func ExampleBasicSenzingSchema_RegisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingschema/senzingschema_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &BasicSenzingSchema{
		SenzingSettings: senzingSettings,
	}
	err = senzingSchema.RegisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Print(err)
	}
	// Output:
}

func ExampleBasicSenzingSchema_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingschema/senzingschema_examples_test.go
	ctx := context.TODO()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &BasicSenzingSchema{
		SenzingSettings: senzingSettings,
	}
	err = senzingSchema.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Print(err)
	}
	// Output:
}

func ExampleBasicSenzingSchema_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingschema/senzingschema_examples_test.go
	ctx := context.TODO()
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &BasicSenzingSchema{
		SenzingSettings: senzingSettings,
	}
	senzingSchema.SetObserverOrigin(ctx, "TestObserver")
	// Output:
}

func ExampleBasicSenzingSchema_UnregisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingschema/senzingschema_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &BasicSenzingSchema{
		SenzingSettings: senzingSettings,
	}
	err = senzingSchema.RegisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Print(err)
	}
	err = senzingSchema.UnregisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Print(err)
	}
	// Output:
}
