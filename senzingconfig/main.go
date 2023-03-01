package senzingconfig

import (
	"context"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-observing/observer"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type SenzingConfig interface {
	Initialize(ctx context.Context) error
	RegisterObserver(ctx context.Context, observer observer.Observer) error
	SetLogLevel(ctx context.Context, logLevel logger.Level) error
	UnregisterObserver(ctx context.Context, observer observer.Observer) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-6502xxxx".
const ProductId = 6502

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for sqlfiler implementation.
var IdMessages = map[int]string{
	1:    "Enter Initialize().",
	2:    "Exit  Initialize() returned (%v).",
	3:    "Enter RegisterObserver(%s).",
	4:    "Exit  RegisterObserver(%s) returned (%v).",
	5:    "Enter SetLogLevel(%v).",
	6:    "Exit  SetLogLevel(%v) returned (%v).",
	7:    "Enter UnregisterObserver(%s).",
	8:    "Exit  UnregisterObserver(%s) returned (%v).",
	901:  "Exit  InitializeSenzingConfiguration() returned (%v).",
	1000: "Entry: %+v",
	2001: "SENZING_ENGINE_CONFIGURATION_JSON: %v",
	2002: "No new Senzing configuration created.  One already exists (%d).",
	2003: "Server listening at %v",
	4001: "Call to net.Listen(tcp, %s) failed.",
	5001: "Failed to serve.",
	8001: "Initialize - config exists",
	8002: "Initialize",
	8003: "RegisterObserver",
	8004: "SetLogLevel",
	8005: "UnregisterObserver",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
