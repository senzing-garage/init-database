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

// Log message prefix.
const Prefix = "init-database.senzingconfig."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for sqlfiler implementation.
var IdMessages = map[int]string{
	10:   "Enter " + Prefix + "InitializeSenzing().",
	11:   "Exit  " + Prefix + "InitializeSenzing(); json.Marshal failed; returned (%v).",
	12:   "Exit  " + Prefix + "InitializeSenzing(); senzingConfig.getDependentServices failed; returned (%v).",
	13:   "Exit  " + Prefix + "InitializeSenzing(); g2Configmgr.GetDefaultConfigID failed; returned (%v).",
	14:   "Exit  " + Prefix + "InitializeSenzing(); Senzing configuration already exists; returned (%v).",
	15:   "Exit  " + Prefix + "InitializeSenzing(); g2Config.Create failed; returned (%v).",
	16:   "Exit  " + Prefix + "InitializeSenzing(); senzingConfig.addDatasources failed; returned (%v).",
	17:   "Exit  " + Prefix + "InitializeSenzing(); g2Config.Save failed; returned (%v).",
	18:   "Exit  " + Prefix + "InitializeSenzing(); g2Configmgr.AddConfig failed; returned (%v).",
	19:   "Exit  " + Prefix + "InitializeSenzing(); g2Configmgr.SetDefaultConfigID failed; returned (%v).",
	29:   "Exit  " + Prefix + "InitializeSenzing() returned (%v).",
	30:   "Enter " + Prefix + "RegisterObserver(%s).",
	31:   "Exit  " + Prefix + "RegisterObserver(%s); json.Marshal failed; returned (%v).",
	32:   "Exit  " + Prefix + "RegisterObserver(%s); senzingConfig.observers.RegisterObserver failed; returned (%v).",
	33:   "Exit  " + Prefix + "RegisterObserver(%s); senzingConfig.getDependentServices failed; returned (%v).",
	34:   "Exit  " + Prefix + "RegisterObserver(%s); g2Config.RegisterObserver failed; returned (%v).",
	35:   "Exit  " + Prefix + "RegisterObserver(%s); g2Configmgr.RegisterObserver failed; returned (%v).",
	39:   "Exit  " + Prefix + "RegisterObserver(%s) returned (%v).",
	40:   "Enter " + Prefix + "SetLogLevel(%s).",
	41:   "Exit  " + Prefix + "SetLogLevel(%s); json.Marshal failed; returned (%v).",
	42:   "Exit  " + Prefix + "SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).",
	43:   "Exit  " + Prefix + "SetLogLevel(%s); senzingConfig.getLogger().SetLogLevel failed; returned (%v).",
	44:   "Exit  " + Prefix + "SetLogLevel(%s); senzingConfig.getDependentServices failed; returned (%v).",
	45:   "Exit  " + Prefix + "SetLogLevel(%s); g2Config.SetLogLevel failed; returned (%v).",
	46:   "Exit  " + Prefix + "SetLogLevel(%s); g2Configmgr.SetLogLevel failed; returned (%v).",
	49:   "Exit  " + Prefix + "SetLogLevel(%s) returned (%v).",
	50:   "Enter " + Prefix + "UnregisterObserver(%s).",
	51:   "Exit  " + Prefix + "UnregisterObserver(%s); json.Marshal failed; returned (%v).",
	52:   "Exit  " + Prefix + "UnregisterObserver(%s); g2Config.UnregisterObserver failed; returned (%v).",
	53:   "Exit  " + Prefix + "UnregisterObserver(%s); g2Configmgr.UnregisterObserver failed; returned (%v).",
	54:   "Exit  " + Prefix + "UnregisterObserver(%s); senzingConfig.observers.UnregisterObserver failed; returned (%v).",
	59:   "Exit  " + Prefix + "UnregisterObserver(%s) returned (%v).",
	1001: Prefix + "InitializeSenzing parameters: %+v",
	1002: Prefix + "RegisterObserver parameters: %+v",
	1003: Prefix + "SetLogLevel parameters: %+v",
	1004: Prefix + "UnregisterObserver parameters: %+v",
	2001: "Added Datasource: %s",
	2002: "No new Senzing configuration created.  One already exists (%d).",
	2003: "Created Senzing configuration: %d named: %s",
	8001: Prefix + "InitializeSenzing - config exists",
	8002: Prefix + "InitializeSenzing",
	8003: Prefix + "RegisterObserver",
	8004: Prefix + "SetLogLevel",
	8005: Prefix + "UnregisterObserver",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
