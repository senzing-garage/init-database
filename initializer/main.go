package initializer

import (
	"context"

	"github.com/senzing/go-observing/observer"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Initializer interface {
	InitializeSenzing(ctx context.Context) error
	RegisterObserver(ctx context.Context, observer observer.Observer) error
	SetLogLevel(ctx context.Context, logLevelName string) error
	UnregisterObserver(ctx context.Context, observer observer.Observer) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-6501xxxx".
const ProductId = 6501

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for g2config implementations.
var IdMessages = map[int]string{
	1:    "Enter Initialize().",
	2:    "Exit  Initialize() returned (%v).",
	3:    "Enter RegisterObserver(%s).",
	4:    "Exit  RegisterObserver(%s) returned (%v).",
	5:    "Enter SetLogLevel(%s).",
	6:    "Exit  SetLogLevel(%s) returned (%v).",
	7:    "Enter UnregisterObserver(%s).",
	8:    "Exit  UnregisterObserver(%s) returned (%v).",
	101:  "Enter InitializeSpecificDatabase().",
	102:  "Exit  InitializeSpecificDatabase() returned (%v).",
	103:  "Exit  InitializeSpecificDatabase() returned (%v).",
	104:  "Exit  InitializeSpecificDatabase() returned (%v).",
	105:  "Exit  InitializeSpecificDatabase() returned (%v).",
	106:  "Exit  InitializeSpecificDatabase() returned (%v).",
	107:  "Enter initializeSpecificDatabaseSqlite(%v).",
	108:  "Exit  initializeSpecificDatabaseSqlite(%v) returned (%v).",
	1000: "Entry: %+v",
	2001: "Created file: %s",
	8001: "Initialize",
	8002: "RegisterObserver",
	8003: "SetLogLevel",
	8004: "UnregisterObserver",
	8005: "InitializeFiles",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
