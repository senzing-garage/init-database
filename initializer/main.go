package initializer

import (
	"context"
	"errors"

	"github.com/senzing-garage/go-observing/observer"
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
const ComponentID = 6501

// Log message prefix.
const Prefix = "init-database.initializer."

// Default gRPC Observer port.
const DefaultGrpcObserverPort = "8260"

const (
	OptionCallerSkip4 = 4
	OptionCallerSkip5 = 5
)

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for szconfig implementations.
var IDMessages = map[int]string{
	10:   "Enter " + Prefix + "Initialize().",
	11:   "Exit  " + Prefix + "Initialize(); json.Marshal failed; returned (%v).",
	12:   "Exit  " + Prefix + "Initialize(); initializerImpl.InitializeSpecificDatabase failed; returned (%v).",
	13:   "Exit  " + Prefix + "Initialize(); senzingSchema.SetLogLevel failed; returned (%v).",
	14:   "Exit  " + Prefix + "Initialize(); senzingSchema.InitializeSenzing failed; returned (%v).",
	15:   "Exit  " + Prefix + "Initialize(); senzingConfig.SetLogLevel failed; returned (%v).",
	16:   "Exit  " + Prefix + "Initialize(); senzingConfig.InitializeSenzing; returned (%v).",
	17:   "Exit  " + Prefix + "Initialize(); initializerImpl.observers.RegisterObserver; returned (%v).",
	18:   "Exit  " + Prefix + "Initialize(); initializerImpl.createGrpcObserver; returned (%v).",
	19:   "Exit  " + Prefix + "Initialize(); initializerImpl.registerObserverSenzingSchema; returned (%v).",
	20:   "Exit  " + Prefix + "Initialize(); initializerImpl.registerObserverSenzingConfig; returned (%v).",
	21:   "Exit  " + Prefix + "Initialize(); os.Stat failed; returned (%v).",
	29:   "Exit  " + Prefix + "Initialize() returned (%v).",
	40:   "Enter " + Prefix + "InitializeSpecificDatabase().",
	41:   "Exit  " + Prefix + "InitializeSpecificDatabase(); json.Marshal failed; returned (%v).",
	42:   "Exit  " + Prefix + "InitializeSpecificDatabase(); settingsparser.New failed; returned (%v).",
	43:   "Exit  " + Prefix + "InitializeSpecificDatabase(); parser.GetDatabaseUrls failed; returned (%v).",
	44:   "Exit  " + Prefix + "InitializeSpecificDatabase(); url.Parse failed; returned (%v).",
	45:   "Exit  " + Prefix + "InitializeSpecificDatabase(); initializerImpl.initializeSpecificDatabaseSqlite; returned (%v).",
	49:   "Exit  " + Prefix + "InitializeSpecificDatabase() returned (%v).",
	50:   "Enter " + Prefix + "RegisterObserver(%s).",
	51:   "Exit  " + Prefix + "RegisterObserver(%s); json.Marshal failed; returned (%v).",
	52:   "Exit  " + Prefix + "RegisterObserver(%s); initializerImpl.observers.RegisterObserver failed; returned (%v).",
	53:   "Exit  " + Prefix + "RegisterObserver(%s); initializerImpl.getSenzingConfig().RegisterObserver failed; returned (%v).",
	54:   "Exit  " + Prefix + "RegisterObserver(%s); initializerImpl.getSenzingSchema().RegisterObserver; returned (%v).",
	59:   "Exit  " + Prefix + "RegisterObserver(%s) returned (%v).",
	60:   "Enter " + Prefix + "SetLogLevel(%s).",
	61:   "Exit  " + Prefix + "SetLogLevel(%s); json.Marshal failed; returned (%v).",
	62:   "Exit  " + Prefix + "SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).",
	63:   "Exit  " + Prefix + "SetLogLevel(%s); initializerImpl.getLogger().SetLogLevel failed; returned (%v).",
	64:   "Exit  " + Prefix + "SetLogLevel(%s); initializerImpl.senzingConfigSingleton.SetLogLevel failed; returned (%v).",
	65:   "Exit  " + Prefix + "SetLogLevel(%s); initializerImpl.getSenzingSchema().SetLogLevel failed; returned (%v).",
	69:   "Exit  " + Prefix + "SetLogLevel(%s) returned (%v).",
	70:   "Enter " + Prefix + "UnregisterObserver(%s).",
	71:   "Exit  " + Prefix + "UnregisterObserver(%s); json.Marshal failed; returned (%v).",
	72:   "Exit  " + Prefix + "UnregisterObserver(%s); initializerImpl.getSenzingConfig().UnregisterObserver failed; returned (%v).",
	73:   "Exit  " + Prefix + "UnregisterObserver(%s); initializerImpl.getSenzingSchema().UnregisterObserver failed; returned (%v).",
	74:   "Exit  " + Prefix + "UnregisterObserver(%s); initializerImpl.observers.UnregisterObserver failed; returned (%v).",
	79:   "Exit  " + Prefix + "UnregisterObserver(%s) returned (%v).",
	80:   "Enter " + Prefix + "SetObserverOrigin(%s).",
	81:   "Exit  " + Prefix + "SetObserverOrigin(%s); json.Marshal failed; returned (%v).",
	89:   "Exit  " + Prefix + "SetObserverOrigin(%s).",
	100:  "Enter " + Prefix + "initializeSpecificDatabaseSqlite(%v).",
	101:  "Exit  " + Prefix + "initializeSpecificDatabaseSqlite(%v); os.Stat failed; returned (%v).",
	102:  "Exit  " + Prefix + "initializeSpecificDatabaseSqlite(%v); os.MkdirAll failed; returned (%v).",
	103:  "Exit  " + Prefix + "initializeSpecificDatabaseSqlite(%v); os.Create failed; returned (%v).",
	109:  "Exit  " + Prefix + "initializeSpecificDatabaseSqlite(%v) returned (%v).",
	1000: Prefix + "Initialize parameters: %+v",
	1001: Prefix + "InitializeSpecificDatabase parameters: %+v",
	1002: Prefix + "RegisterObserver parameters: %+v",
	1003: Prefix + "SetLogLevel parameters: %+v",
	1004: Prefix + "SetObserverOrigin parameters: %+v",
	1005: Prefix + "UnregisterObserver parameters: %+v",
	1011: Prefix + "Initialize(); json.Marshal failed; Error: %v.",
	1012: Prefix + "Initialize(); initializerImpl.InitializeSpecificDatabase failed; Error: %v.",
	1013: Prefix + "Initialize(); initializerImpl.getSenzingSchema failed; Error: %v.",
	1014: Prefix + "Initialize(); senzingSchema.InitializeSenzing failed; Error: %v.",
	1015: Prefix + "Initialize(); initializerImpl.getSenzingConfig failed; Error: %v.",
	1016: Prefix + "Initialize(); senzingConfig.InitializeSenzing; Error: %v.",
	1017: Prefix + "Initialize(); initializerImpl.observers.RegisterObserver; returned (%v).",
	1018: Prefix + "Initialize(); initializerImpl.createGrpcObserver; returned (%v).",
	1041: Prefix + "InitializeSpecificDatabase(); json.Marshal failed; Error: %v.",
	1042: Prefix + "InitializeSpecificDatabase(); settingsparser.New failed; Error: %v.",
	1043: Prefix + "InitializeSpecificDatabase(); parser.GetDatabaseUrls failed; Error: %v.",
	1044: Prefix + "InitializeSpecificDatabase(); url.Parse failed; Error: %v.",
	1045: Prefix + "InitializeSpecificDatabase(); initializerImpl.initializeSpecificDatabaseSqlite; Error: %v.",
	1051: Prefix + "RegisterObserver(%s); json.Marshal failed; Error: %v.",
	1052: Prefix + "RegisterObserver(%s); initializerImpl.observers.RegisterObserver failed; Error: %v.",
	1053: Prefix + "RegisterObserver(%s); initializerImpl.getSenzingConfig().RegisterObserver failed; Error: %v.",
	1054: Prefix + "RegisterObserver(%s); initializerImpl.getSenzingSchema().RegisterObserver; Error: %v.",
	1061: Prefix + "SetLogLevel(%s); json.Marshal failed; Error: %v.",
	1062: Prefix + "SetLogLevel(%s); logging.IsValidLogLevelName failed; Error: %v.",
	1063: Prefix + "SetLogLevel(%s); initializerImpl.getLogger().SetLogLevel failed; Error: %v.",
	1064: Prefix + "SetLogLevel(%s); initializerImpl.senzingConfigSingleton.SetLogLevel failed; Error: %v.",
	1065: Prefix + "SetLogLevel(%s); initializerImpl.getSenzingSchema().SetLogLevel failed; Error: %v.",
	1071: Prefix + "UnregisterObserver(%s); json.Marshal failed; Error: %v.",
	1072: Prefix + "UnregisterObserver(%s); initializerImpl.getSenzingConfig().UnregisterObserver failed; Error: %v.",
	1073: Prefix + "UnregisterObserver(%s); initializerImpl.getSenzingSchema().UnregisterObserver failed; Error: %v.",
	1074: Prefix + "UnregisterObserver(%s); initializerImpl.observers.UnregisterObserver failed; Error: %v.",
	1075: Prefix + "Initialize(); os.Stat failed; Error: %v.",
	1081: Prefix + "SetObserverOrigin(%s); json.Marshal failed; Error: %v.",
	1101: Prefix + "initializeSpecificDatabaseSqlite(%v); os.Stat failed; returned (%v).",
	1102: Prefix + "initializeSpecificDatabaseSqlite(%v); os.MkdirAll failed; returned (%v).",
	1103: Prefix + "initializeSpecificDatabaseSqlite(%v); os.Create failed; returned (%v).",
	2001: "Created file: %s",
	3001: "SQL file does not exist: %s",
	8001: Prefix + "Initialize Observer URL",
	8002: Prefix + "Initialize",
	8003: Prefix + "RegisterObserver",
	8004: Prefix + "SetLogLevel",
	8005: Prefix + "SetObserverOrigin",
	8006: Prefix + "UnregisterObserver",
	8010: Prefix + "initializeSpecificDatabaseSqlite",
}

// Status strings for specific messages.
var IDStatuses = map[int]string{}

var errForPackage = errors.New("initializer")
