//go:build linux

package senzingschema

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSenzingSchemaImpl_InitializeSenzing() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingschema/senzingschema_examples_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &SenzingSchemaImpl{
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
	// Output:
}

func ExampleSenzingSchemaImpl_RegisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingschema/senzingschema_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingSchema.RegisterObserver(ctx, anObserver)
	if err != nil {
		fmt.Print(err)
	}
	// Output:
}

func ExampleSenzingSchemaImpl_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingschema/senzingschema_examples_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingSchema.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Print(err)
	}
	// Output:
}

func ExampleSenzingSchemaImpl_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingschema/senzingschema_examples_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	senzingSchema.SetObserverOrigin(ctx, "TestObserver")
	// Output:
}

func ExampleSenzingSchemaImpl_UnregisterObserver() {
	// For more information, visit https://github.com/senzing-garage/init-database/blob/main/senzingschema/senzingschema_examples_test.go
	ctx := context.TODO()
	anObserver := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
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
