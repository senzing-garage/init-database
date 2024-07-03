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

func ExampleInitializerImpl_Initialize() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/initializer/initializer_examples_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := settings.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	anInitializer := &InitializerImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
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

func ExampleInitializerImpl_RegisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/initializer/initializer_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := settings.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	anInitializer := &InitializerImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = anInitializer.RegisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Print(err)
	}
	// Output:
}

func ExampleInitializerImpl_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/initializer/initializer_examples_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := settings.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	anInitializer := &InitializerImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = anInitializer.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleInitializerImpl_SetObserverOrigin() {
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := settings.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	anInitializer := &InitializerImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	anInitializer.SetObserverOrigin(ctx, "TestObserver")
	// Output:
}

func ExampleInitializerImpl_UnregisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/initializer/initializer_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := settings.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	anInitializer := &InitializerImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = anInitializer.RegisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Print(err)
	}
	err = anInitializer.UnregisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Print(err)
	}
	// Output:
}
