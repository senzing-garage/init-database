package senzingschema

import (
	"context"

	"github.com/senzing/go-observing/observer"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type SenzingSchema interface {
	InitializeSenzing(ctx context.Context) error
	RegisterObserver(ctx context.Context, observer observer.Observer) error
	SetLogLevel(ctx context.Context, logLevelName string) error
	UnregisterObserver(ctx context.Context, observer observer.Observer) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-6503xxxx".
const ProductId = 6503

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for g2config implementations.
var IdMessages = map[int]string{
	1:    "Enter InitializeSenzing().",
	2:    "Exit  InitializeSenzing() returned (%v).",
	3:    "Enter RegisterObserver(%s).",
	4:    "Exit  RegisterObserver(%s) returned (%v).",
	5:    "Enter SetLogLevel(%s).",
	6:    "Exit  SetLogLevel(%s) returned (%v).",
	7:    "Enter UnregisterObserver(%s).",
	8:    "Exit  UnregisterObserver(%s) returned (%v).",
	901:  "Exit  InitializeSenzingConfiguration() returned (%v).",
	1000: "Entry: %+v",
	2001: "SENZING_ENGINE_CONFIGURATION_JSON: %v",
	2002: "Sent SQL in %s to database %s",
	8001: "InitializeSenzing",
	8002: "RegisterObserver",
	8003: "SetLogLevel",
	8004: "UnregisterObserver",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
