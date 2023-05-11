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
const ComponentId = 6503

// Log message prefix.
const Prefix = "init-database.senzingconfig."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for g2config implementations.
var IdMessages = map[int]string{
	10:   "Enter " + Prefix + "InitializeSenzing().",
	11:   "Exit  " + Prefix + "InitializeSenzing(); json.Marshal failed; returned (%v).",
	12:   "Exit  " + Prefix + "InitializeSenzing(); engineconfigurationjsonparser.New failed; returned (%v).",
	13:   "Exit  " + Prefix + "InitializeSenzing(); parser.GetResourcePath failed; returned (%v).",
	14:   "Exit  " + Prefix + "InitializeSenzing(); parser.GetDatabaseUrls failed; returned (%v).",
	15:   "Exit  " + Prefix + "InitializeSenzing(); senzingSchema.processDatabase failed; returned (%v).",
	19:   "Exit  " + Prefix + "InitializeSenzing() returned (%v).",
	20:   "Enter " + Prefix + "RegisterObserver(%s).",
	21:   "Exit  " + Prefix + "RegisterObserver(%s); json.Marshal failed; returned (%v).",
	22:   "Exit  " + Prefix + "RegisterObserver(%s); senzingSchema.observers.RegisterObserver failed; returned (%v).",
	29:   "Exit  " + Prefix + "RegisterObserver(%s) returned (%v).",
	30:   "Enter " + Prefix + "SetLogLevel(%s).",
	31:   "Exit  " + Prefix + "SetLogLevel(%s); json.Marshal failed; returned (%v).",
	32:   "Exit  " + Prefix + "SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).",
	39:   "Exit  " + Prefix + "SetLogLevel(%s) returned (%v).",
	40:   "Enter " + Prefix + "UnregisterObserver(%s).",
	41:   "Exit  " + Prefix + "UnregisterObserver(%s); json.Marshal failed; returned (%v).",
	42:   "Exit  " + Prefix + "UnregisterObserver(%s); senzingSchema.observers.UnregisterObserver failed; returned (%v).",
	49:   "Exit  " + Prefix + "UnregisterObserver(%s) returned (%v).",
	100:  "Enter " + Prefix + "processDatabase(%s, %s).",
	101:  "Exit  " + Prefix + "processDatabase(%s, %s); url.Parse failed; returned (%v).",
	102:  "Exit  " + Prefix + "processDatabase(%s, %s); connector.NewConnector failed; returned (%v).",
	103:  "Exit  " + Prefix + "processDatabase(%s, %s); sqlExecutor.SetLogLevel failed; returned (%v).",
	104:  "Exit  " + Prefix + "processDatabase(%s, %s); sqlExecutor.RegisterObserver failed; returned (%v).",
	105:  "Exit  " + Prefix + "processDatabase(%s, %s); sqlExecutor.ProcessFileName failed; returned (%v).",
	109:  "Exit  " + Prefix + "processDatabase(%s, %s) returned (%v).",
	1001: Prefix + "InitializeSenzing parameters: %+v",
	1002: Prefix + "RegisterObserver parameters: %+v",
	1003: Prefix + "SetLogLevel parameters: %+v",
	1004: Prefix + "UnregisterObserver parameters: %+v",
	1011: Prefix + "InitializeSenzing(); json.Marshal failed; returned (%v).",
	1012: Prefix + "InitializeSenzing(); engineconfigurationjsonparser.New failed; returned (%v).",
	1013: Prefix + "InitializeSenzing(); parser.GetResourcePath failed; returned (%v).",
	1014: Prefix + "InitializeSenzing(); parser.GetDatabaseUrls failed; returned (%v).",
	1015: Prefix + "InitializeSenzing(); senzingSchema.processDatabase failed; returned (%v).",
	1021: Prefix + "RegisterObserver(%s); json.Marshal failed; returned (%v).",
	1022: Prefix + "RegisterObserver(%s); senzingSchema.observers.RegisterObserver failed; returned (%v).",
	1031: Prefix + "SetLogLevel(%s); json.Marshal failed; returned (%v).",
	1032: Prefix + "SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).",
	1041: Prefix + "UnregisterObserver(%s); json.Marshal failed; returned (%v).",
	1042: Prefix + "UnregisterObserver(%s); senzingSchema.observers.UnregisterObserver failed; returned (%v).",
	2001: "Sent SQL in %s to database %s",
	8001: Prefix + "InitializeSenzing",
	8002: Prefix + "RegisterObserver",
	8003: Prefix + "SetLogLevel",
	8004: Prefix + "UnregisterObserver",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
