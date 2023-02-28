package initializer

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-logging/logger"
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

func TestInitializerImpl_Initialize(test *testing.T) {
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	testObject := &InitializerImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	testObject.SetLogLevel(ctx, logger.LevelTrace)
	testObject.Initialize(ctx)
}

func TestInitializerImpl_RegisterObserver(test *testing.T) {
	ctx := context.TODO()

	observer1 := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}

	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	testObject := &InitializerImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	testObject.SetLogLevel(ctx, logger.LevelInfo)
	testObject.RegisterObserver(ctx, observer1)
	testObject.Initialize(ctx)
}

func TestInitializerImpl_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()

	observer1 := &observer.ObserverNull{
		Id:       "Observer 1",
		IsSilent: true,
	}

	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	testObject := &InitializerImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	testObject.SetLogLevel(ctx, logger.LevelInfo)
	testObject.RegisterObserver(ctx, observer1)
	testObject.Initialize(ctx)
	testObject.UnregisterObserver(ctx, observer1)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleInitializerImpl_Initialize() {
	// For more information, visit https://github.com/Senzing/initdatabase/blob/main/initializer/initializer_test.go
	ctx := context.TODO()
	senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Print(err)
	}
	anInitializer := &InitializerImpl{
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
	}
	err = anInitializer.Initialize(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
