package senzingschema

import (
	"context"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-observing/observer"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type SenzingSchema interface {
	Initialize(ctx context.Context) error
	RegisterObserver(ctx context.Context, observer observer.Observer) error
	SetLogLevel(ctx context.Context, logLevel logger.Level) error
	UnregisterObserver(ctx context.Context, observer observer.Observer) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-6212xxxx".
const ProductId = 6212

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for g2config implementations.
var IdMessages = map[int]string{
	1:    "Enter InitializeSenzingConfiguration().",
	2:    "Exit  InitializeSenzingConfiguration() returned (%v).",
	3:    "Enter Initialize().",
	4:    "Exit  Initialize() returned (%v).",
	5:    "Enter RegisterObserver(%s).",
	6:    "Exit  RegisterObserver(%s) returned (%v).",
	7:    "Enter SetLogLevel(%v).",
	8:    "Exit  SetLogLevel(%v) returned (%v).",
	9:    "Enter UnregisterObserver(%s).",
	10:   "Exit  UnregisterObserver(%s) returned (%v).",
	901:  "Exit  InitializeSenzingConfiguration() returned (%v).",
	2000: "Entry: %+v",
	2001: "SENZING_ENGINE_CONFIGURATION_JSON: %v",
	2002: "No new Senzing configuration created.  One already exists (%d).",
	2003: "Server listening at %v",
	4001: "Call to net.Listen(tcp, %s) failed.",
	5001: "Failed to serve.",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
