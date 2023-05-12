package senzingschema

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-observing/observer"
)

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
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	testObject := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	testObject.SetLogLevel(ctx, logging.LevelInfoName)
	testObject.InitializeSenzing(ctx)
}

func TestSenzingSchemaImpl_RegisterObserver(test *testing.T) {
	ctx := context.TODO()

	observer1 := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}

	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	testObject := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	testObject.SetLogLevel(ctx, logging.LevelInfoName)
	testObject.RegisterObserver(ctx, observer1)
	testObject.InitializeSenzing(ctx)
}

func TestSenzingSchemaImpl_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	testObject := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	testObject.SetObserverOrigin(ctx, "TestObserver")
}

func TestSenzingSchemaImpl_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()

	observer1 := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}

	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	testObject := &SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	testObject.SetLogLevel(ctx, logging.LevelInfoName)
	testObject.RegisterObserver(ctx, observer1)
	testObject.InitializeSenzing(ctx)
	testObject.UnregisterObserver(ctx, observer1)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSenzingSchemaImpl_InitializeSenzing() {
	// For more information, visit https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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
	// For more information, visit https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema_test.go
	ctx := context.TODO()
	anObserver := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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
	// For more information, visit https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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
	// For more information, visit https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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
	// For more information, visit https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema_test.go
	ctx := context.TODO()
	anObserver := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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
