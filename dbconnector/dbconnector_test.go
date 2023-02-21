package dbconnector

import (
	"fmt"
	"os"
	"testing"

	"github.com/senzing/go-logging/messagelogger"
	"golang.org/x/net/context"
)

var (
	localLogger messagelogger.MessageLoggerInterface
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
	databaseConnector := &Sqlite{
		Name:     "bob",
		Filename: "/tmp/sqlite/G2C.db",
	}
	databaseConnector.Connect(ctx)
}
