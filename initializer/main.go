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
const Prefix = "init-database.initializer."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for g2config implementations.
var IdMessages = map[int]string{
	10:   "Enter " + Prefix + "Initialize().",
	11:   "Exit  " + Prefix + "Initialize(); json.Marshal failed; returned (%v).",
	12:   "Exit  " + Prefix + "Initialize(); initializerImpl.InitializeSpecificDatabase failed; returned (%v).",
	13:   "Exit  " + Prefix + "Initialize(); senzingSchema.SetLogLevel failed; returned (%v).",
	14:   "Exit  " + Prefix + "Initialize(); senzingSchema.InitializeSenzing failed; returned (%v).",
	15:   "Exit  " + Prefix + "Initialize(); senzingConfig.SetLogLevel failed; returned (%v).",
	16:   "Exit  " + Prefix + "Initialize(); senzingConfig.InitializeSenzing; returned (%v).",
	19:   "Exit  " + Prefix + "Initialize() returned (%v).",
	20:   "Enter " + Prefix + "InitializeSpecificDatabase().",
	21:   "Exit  " + Prefix + "InitializeSpecificDatabase(); json.Marshal failed; returned (%v).",
	22:   "Exit  " + Prefix + "InitializeSpecificDatabase(); engineconfigurationjsonparser.New failed; returned (%v).",
	23:   "Exit  " + Prefix + "InitializeSpecificDatabase(); parser.GetDatabaseUrls failed; returned (%v).",
	24:   "Exit  " + Prefix + "InitializeSpecificDatabase(); url.Parse failed; returned (%v).",
	25:   "Exit  " + Prefix + "InitializeSpecificDatabase(); initializerImpl.initializeSpecificDatabaseSqlite; returned (%v).",
	29:   "Exit  " + Prefix + "InitializeSpecificDatabase() returned (%v).",
	30:   "Enter " + Prefix + "RegisterObserver().",
	31:   "Exit  " + Prefix + "RegisterObserver(); json.Marshal failed; returned (%v).",
	32:   "Exit  " + Prefix + "RegisterObserver(); initializerImpl.observers.RegisterObserver failed; returned (%v).",
	33:   "Exit  " + Prefix + "RegisterObserver(); initializerImpl.getSenzingConfig().RegisterObserver failed; returned (%v).",
	34:   "Exit  " + Prefix + "RegisterObserver(); initializerImpl.getSenzingSchema().RegisterObserver; returned (%v).",
	39:   "Exit  " + Prefix + "RegisterObserver() returned (%v).",
	40:   "Enter " + Prefix + "SetLogLevel().",
	41:   "Exit  " + Prefix + "SetLogLevel(); json.Marshal failed; returned (%v).",
	42:   "Exit  " + Prefix + "SetLogLevel(); logging.IsValidLogLevelName failed; returned (%v).",
	43:   "Exit  " + Prefix + "SetLogLevel(); initializerImpl.getLogger().SetLogLevel failed; returned (%v).",
	44:   "Exit  " + Prefix + "SetLogLevel(); initializerImpl.senzingConfigSingleton.SetLogLevel failed; returned (%v).",
	45:   "Exit  " + Prefix + "SetLogLevel(); initializerImpl.getSenzingSchema().SetLogLevel failed; returned (%v).",
	49:   "Exit  " + Prefix + "SetLogLevel() returned (%v).",
	50:   "Enter " + Prefix + "UnregisterObserver().",
	51:   "Exit  " + Prefix + "UnregisterObserver(); json.Marshal failed; returned (%v).",
	52:   "Exit  " + Prefix + "UnregisterObserver(); initializerImpl.getSenzingConfig().UnregisterObserver failed; returned (%v).",
	53:   "Exit  " + Prefix + "UnregisterObserver(); initializerImpl.getSenzingSchema().UnregisterObserver failed; returned (%v).",
	54:   "Exit  " + Prefix + "UnregisterObserver(); initializerImpl.observers.UnregisterObserver failed; returned (%v).",
	59:   "Exit  " + Prefix + "UnregisterObserver() returned (%v).",
	1000: Prefix + "Initialize parameters: %+v",
	1001: Prefix + "InitializeSpecificDatabase parameters: %+v",
	1002: Prefix + "RegisterObserver parameters: %+v",
	1003: Prefix + "SetLogLevel parameters: %+v",
	1004: Prefix + "UnregisterObserver parameters: %+v",
	1011: Prefix + "Initialize(); json.Marshal failed; Error: %v.",
	1012: Prefix + "Initialize(); initializerImpl.InitializeSpecificDatabase failed; Error: %v.",
	1013: Prefix + "Initialize(); senzingSchema.SetLogLevel failed; Error: %v.",
	1014: Prefix + "Initialize(); senzingSchema.InitializeSenzing failed; Error: %v.",
	1015: Prefix + "Initialize(); senzingConfig.SetLogLevel failed; Error: %v.",
	1016: Prefix + "Initialize(); senzingConfig.InitializeSenzing; Error: %v.",
	1021: Prefix + "InitializeSpecificDatabase(); json.Marshal failed; Error: %v.",
	1022: Prefix + "InitializeSpecificDatabase(); engineconfigurationjsonparser.New failed; Error: %v.",
	1023: Prefix + "InitializeSpecificDatabase(); parser.GetDatabaseUrls failed; Error: %v.",
	1024: Prefix + "InitializeSpecificDatabase(); url.Parse failed; Error: %v.",
	1025: Prefix + "InitializeSpecificDatabase(); initializerImpl.initializeSpecificDatabaseSqlite; Error: %v.",
	1031: Prefix + "RegisterObserver(); json.Marshal failed; Error: %v.",
	1032: Prefix + "RegisterObserver(); initializerImpl.observers.RegisterObserver failed; Error: %v.",
	1033: Prefix + "RegisterObserver(); initializerImpl.getSenzingConfig().RegisterObserver failed; Error: %v.",
	1034: Prefix + "RegisterObserver(); initializerImpl.getSenzingSchema().RegisterObserver; Error: %v.",
	1041: Prefix + "SetLogLevel(); json.Marshal failed; Error: %v.",
	1042: Prefix + "SetLogLevel(); logging.IsValidLogLevelName failed; Error: %v.",
	1043: Prefix + "SetLogLevel(); initializerImpl.getLogger().SetLogLevel failed; Error: %v.",
	1044: Prefix + "SetLogLevel(); initializerImpl.senzingConfigSingleton.SetLogLevel failed; Error: %v.",
	1045: Prefix + "SetLogLevel(); initializerImpl.getSenzingSchema().SetLogLevel failed; Error: %v.",
	1051: Prefix + "UnregisterObserver(); json.Marshal failed; Error: %v.",
	1052: Prefix + "UnregisterObserver(); initializerImpl.getSenzingConfig().UnregisterObserver failed; Error: %v.",
	1053: Prefix + "UnregisterObserver(); initializerImpl.getSenzingSchema().UnregisterObserver failed; Error: %v.",
	1054: Prefix + "UnregisterObserver(); initializerImpl.observers.UnregisterObserver failed; Error: %v.",
	2001: "Created file: %s",
	8001: Prefix + "Initialize",
	8002: Prefix + "RegisterObserver",
	8003: Prefix + "SetLogLevel",
	8004: Prefix + "UnregisterObserver",
	8005: Prefix + "InitializeFiles",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
