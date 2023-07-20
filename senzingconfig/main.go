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
	SetObserverOrigin(ctx context.Context, origin string)
	UnregisterObserver(ctx context.Context, observer observer.Observer) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-6502xxxx".
const ComponentId = 6502

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
	20:   "Exit  " + Prefix + "InitializeSenzing(); engineconfigurationjsonparser.New failed; returned (%v).",
	21:   "Exit  " + Prefix + "InitializeSenzing(); engineconfigurationjsonparser.GetResourcePath failed; returned (%v).",
	22:   "Exit  " + Prefix + "InitializeSenzing(); os.Stat failed; returned (%v).",
	23:   "Exit  " + Prefix + "InitializeSenzing(); copyFile when backing up failed; returned (%v).",
	24:   "Exit  " + Prefix + "InitializeSenzing(); copyFile when replacing template/g2config.json failed; returned (%v).",
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
	60:   "Enter " + Prefix + "SetObserverOrigin(%s).",
	61:   "Exit  " + Prefix + "SetObserverOrigin(%s); json.Marshal failed; returned (%v).",
	69:   "Exit  " + Prefix + "SetObserverOrigin(%s).",
	1001: Prefix + "InitializeSenzing parameters: %+v",
	1002: Prefix + "RegisterObserver parameters: %+v",
	1003: Prefix + "SetLogLevel parameters: %+v",
	1004: Prefix + "SetObserverOrigin parameters: %+v",
	1005: Prefix + "UnregisterObserver parameters: %+v",
	1011: Prefix + "Initialize(); json.Marshal failed; Error: %v.",
	1012: Prefix + "Initialize(); senzingConfig.getDependentServices failed; Error: %v.",
	1013: Prefix + "Initialize(); g2Configmgr.GetDefaultConfigID failed; Error: %v.",
	1014: Prefix + "Initialize(); senzingSchema.InitializeSenzing failed; Error: %v.",
	1015: Prefix + "Initialize(); g2Config.Create failed; Error: %v.",
	1016: Prefix + "Initialize(); senzingConfig.addDatasources failed; Error: %v.",
	1017: Prefix + "Initialize(); g2Config.Save failed; Error: %v.",
	1018: Prefix + "Initialize(); g2Configmgr.AddConfig failed; Error: %v.",
	1019: Prefix + "Initialize(); g2Configmgr.SetDefaultConfigID failed; Error: %v.",
	1020: Prefix + "Initialize(); engineconfigurationjsonparser.New failed; Error: %v.",
	1021: Prefix + "Initialize(); engineconfigurationjsonparser.GetResourcePath failed; Error: %v.",
	1022: Prefix + "Initialize(); os.Stat failed; Error: %v.",
	1023: Prefix + "Initialize(); copyFile when backing up failed; Error: %v.",
	1024: Prefix + "Initialize(); copyFile when replacing template/g2config.json failed; Error: %v.",
	1031: Prefix + "RegisterObserver(%s); json.Marshal failed; returned (%v).",
	1032: Prefix + "RegisterObserver(%s); senzingConfig.observers.RegisterObserver failed; returned (%v).",
	1033: Prefix + "RegisterObserver(%s); senzingConfig.getDependentServices failed; returned (%v).",
	1034: Prefix + "RegisterObserver(%s); g2Config.RegisterObserver failed; returned (%v).",
	1035: Prefix + "RegisterObserver(%s); g2Configmgr.RegisterObserver failed; returned (%v).",
	1041: Prefix + "SetLogLevel(%s); json.Marshal failed; returned (%v).",
	1042: Prefix + "SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).",
	1043: Prefix + "SetLogLevel(%s); senzingConfig.getLogger().SetLogLevel failed; returned (%v).",
	1044: Prefix + "SetLogLevel(%s); senzingConfig.getDependentServices failed; returned (%v).",
	1045: Prefix + "SetLogLevel(%s); g2Config.SetLogLevel failed; returned (%v).",
	1046: Prefix + "SetLogLevel(%s); g2Configmgr.SetLogLevel failed; returned (%v).",
	1051: Prefix + "UnregisterObserver(%s); json.Marshal failed; returned (%v).",
	1052: Prefix + "UnregisterObserver(%s); g2Config.UnregisterObserver failed; returned (%v).",
	1053: Prefix + "UnregisterObserver(%s); g2Configmgr.UnregisterObserver failed; returned (%v).",
	1054: Prefix + "UnregisterObserver(%s); senzingConfig.observers.UnregisterObserver failed; returned (%v).",
	2001: "Added Datasource: %s",
	2002: "No new Senzing configuration created.  One already exists (%d).",
	2003: "Created Senzing configuration: %d named: %s",
	2004: "Copied file %s to %s",
	2005: "%s and %s have same content.  No file manipulation needed.",
	4001: "When comparing %s and %s, an error occurred. Assuming files not equal.",
	5001: "File does not exist: %s [SENZING_TOOLS_ENGINE_CONFIGURATION_FILE]",
	5002: "Could not backup %s to %s",
	5003: "Could not copy %s to %s [SENZING_TOOLS_ENGINE_CONFIGURATION_FILE]",
	8001: Prefix + "InitializeSenzing - config exists",
	8002: Prefix + "InitializeSenzing",
	8003: Prefix + "RegisterObserver",
	8004: Prefix + "SetLogLevel",
	8005: Prefix + "SetObserverOrigin",
	8006: Prefix + "UnregisterObserver",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
