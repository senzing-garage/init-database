package sqlfiler

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/initdatabase/dbconnector"
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

	databaseConnector := &dbconnector.Sqlite{
		Name:     "bob",
		Filename: "/tmp/sqlite/G2C.db",
	}

	testObject := &SqlfilerImpl{
		LogLevel:          logger.LevelTrace,
		DatabaseConnector: databaseConnector,
	}
	testObject.ProcessFileName(ctx, "/opt/senzing/g2/resources/schema/g2core-schema-sqlite-create.sql")
}
