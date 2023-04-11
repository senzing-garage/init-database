package senzingconfig

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/init-database/senzingschema"
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
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	senzingSchema := &senzingschema.SenzingSchemaImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingSchema.SetLogLevel(ctx, logger.LevelInfo)
	if err != nil {
		fmt.Println(err)
	}
	err = senzingSchema.Initialize(ctx)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSenzingConfigImpl_Initialize_withDatasources() {
	// For more information, visit https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		DataSources:                    []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"},
	}
	err = senzingConfig.SetLogLevel(ctx, logger.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	err = senzingConfig.Initialize(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSenzingConfigImpl_Initialize() {
	// For more information, visit https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingConfig.SetLogLevel(ctx, logger.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	err = senzingConfig.Initialize(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSenzingConfigImpl_RegisterObserver() {
	// For more information, visit https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig_test.go
	ctx := context.TODO()
	anObserver := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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
	// For more information, visit https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	senzingConfig := &SenzingConfigImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = senzingConfig.SetLogLevel(ctx, logger.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSenzingConfigImpl_UnregisterObserver() {
	// For more information, visit https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig_test.go
	ctx := context.TODO()
	anObserver := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
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
