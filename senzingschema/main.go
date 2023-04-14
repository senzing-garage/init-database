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
	10:   "Enter InitializeSenzing().",
	11:   "Exit  InitializeSenzing(); json.Marshal failed; returned (%v).",
	12:   "Exit  InitializeSenzing(); engineconfigurationjsonparser.New failed; returned (%v).",
	13:   "Exit  InitializeSenzing(); parser.GetResourcePath failed; returned (%v).",
	14:   "Exit  InitializeSenzing(); parser.GetDatabaseUrls failed; returned (%v).",
	15:   "Exit  InitializeSenzing(); senzingSchema.processDatabase failed; returned (%v).",
	19:   "Exit  InitializeSenzing() returned (%v).",
	20:   "Enter RegisterObserver(%s).",
	21:   "Exit  RegisterObserver(%s); json.Marshal failed; returned (%v).",
	22:   "Exit  RegisterObserver(%s); senzingSchema.observers.RegisterObserver failed; returned (%v).",
	29:   "Exit  RegisterObserver(%s) returned (%v).",
	30:   "Enter SetLogLevel(%s).",
	31:   "Exit  SetLogLevel(%s); json.Marshal failed; returned (%v).",
	32:   "Exit  SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).",
	39:   "Exit  SetLogLevel(%s) returned (%v).",
	40:   "Enter UnregisterObserver(%s).",
	41:   "Exit  UnregisterObserver(%s); json.Marshal failed; returned (%v).",
	42:   "Exit  UnregisterObserver(%s); senzingSchema.observers.UnregisterObserver failed; returned (%v).",
	49:   "Exit  UnregisterObserver(%s) returned (%v).",
	100:  "Enter processDatabase().",
	101:  "Exit  InitializeSenzing(); url.Parse failed; returned (%v).",
	102:  "Exit  InitializeSenzing(); connector.NewConnector failed; returned (%v).",
	103:  "Exit  InitializeSenzing(); sqlExecutor.SetLogLevel failed; returned (%v).",
	104:  "Exit  InitializeSenzing(); sqlExecutor.RegisterObserver failed; returned (%v).",
	105:  "Exit  InitializeSenzing(); sqlExecutor.ProcessFileName failed; returned (%v).",
	109:  "Exit  InitializeSenzing() returned (%v).",
	1001: "InitializeSenzing parameters: %+v",
	1002: "RegisterObserver parameters: %+v",
	1003: "SetLogLevel parameters: %+v",
	1004: "UnregisterObserver parameters: %+v",
	2001: "Sent SQL in %s to database %s",
	8001: "InitializeSenzing",
	8002: "RegisterObserver",
	8003: "SetLogLevel",
	8004: "UnregisterObserver",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
