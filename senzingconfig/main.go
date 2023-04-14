package senzingconfig

import (
	"context"

	"github.com/senzing/go-observing/observer"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type SenzingConfig interface {
	InitializeSenzing(ctx context.Context) error
	RegisterObserver(ctx context.Context, observer observer.Observer) error
	SetLogLevel(ctx context.Context, logLevelName string) error
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
	10:   "Enter InitializeSenzing().",
	11:   "Exit  InitializeSenzing(); json.Marshal failed; returned (%v).",
	12:   "Exit  InitializeSenzing(); senzingConfig.getDependentServices failed; returned (%v).",
	13:   "Exit  InitializeSenzing(); g2Configmgr.GetDefaultConfigID failed; returned (%v).",
	14:   "Exit  InitializeSenzing(); Senzing configuration already exists; returned (%v).",
	15:   "Exit  InitializeSenzing(); g2Config.Create failed; returned (%v).",
	16:   "Exit  InitializeSenzing(); senzingConfig.addDatasources failed; returned (%v).",
	17:   "Exit  InitializeSenzing(); g2Config.Save failed; returned (%v).",
	18:   "Exit  InitializeSenzing(); g2Configmgr.AddConfig failed; returned (%v).",
	19:   "Exit  InitializeSenzing(); g2Configmgr.SetDefaultConfigID failed; returned (%v).",
	29:   "Exit  InitializeSenzing() returned (%v).",
	30:   "Enter RegisterObserver(%s).",
	31:   "Exit  RegisterObserver(%s); json.Marshal failed; returned (%v).",
	32:   "Exit  RegisterObserver(%s); senzingConfig.observers.RegisterObserver failed; returned (%v).",
	33:   "Exit  RegisterObserver(%s); senzingConfig.getDependentServices failed; returned (%v).",
	34:   "Exit  RegisterObserver(%s); g2Config.RegisterObserver failed; returned (%v).",
	35:   "Exit  RegisterObserver(%s); g2Configmgr.RegisterObserver failed; returned (%v).",
	39:   "Exit  RegisterObserver(%s) returned (%v).",
	40:   "Enter SetLogLevel(%s).",
	41:   "Exit  SetLogLevel(%s); json.Marshal failed; returned (%v).",
	42:   "Exit  SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).",
	43:   "Exit  SetLogLevel(%s); senzingConfig.getLogger().SetLogLevel failed; returned (%v).",
	44:   "Exit  SetLogLevel(%s); senzingConfig.getDependentServices failed; returned (%v).",
	45:   "Exit  SetLogLevel(%s); g2Config.SetLogLevel failed; returned (%v).",
	46:   "Exit  SetLogLevel(%s); g2Configmgr.SetLogLevel failed; returned (%v).",
	49:   "Exit  SetLogLevel(%s) returned (%v).",
	50:   "Enter UnregisterObserver(%s).",
	51:   "Exit  UnregisterObserver(%s); json.Marshal failed; returned (%v).",
	52:   "Exit  UnregisterObserver(%s); g2Config.UnregisterObserver failed; returned (%v).",
	53:   "Exit  UnregisterObserver(%s); g2Configmgr.UnregisterObserver failed; returned (%v).",
	54:   "Exit  UnregisterObserver(%s); senzingConfig.observers.UnregisterObserver failed; returned (%v).",
	59:   "Exit  UnregisterObserver(%s) returned (%v).",
	1001: "InitializeSenzing parameters: %+v",
	1002: "RegisterObserver parameters: %+v",
	1003: "SetLogLevel parameters: %+v",
	1004: "UnregisterObserver parameters: %+v",
	2001: "Added Datasource: %s",
	2002: "No new Senzing configuration created.  One already exists (%d).",
	2003: "Created Senzing configuration: %d named: %s",
	8001: "InitializeSenzing - config exists",
	8002: "InitializeSenzing",
	8003: "RegisterObserver",
	8004: "SetLogLevel",
	8005: "UnregisterObserver",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
