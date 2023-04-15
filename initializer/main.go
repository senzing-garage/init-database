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

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for g2config implementations.
var IdMessages = map[int]string{
	10:   "Enter Initialize().",
	11:   "Exit  Initialize(); json.Marshal failed; returned (%v).",
	12:   "Exit  Initialize(); initializerImpl.InitializeSpecificDatabase failed; returned (%v).",
	13:   "Exit  Initialize(); senzingSchema.SetLogLevel failed; returned (%v).",
	14:   "Exit  Initialize(); senzingSchema.InitializeSenzing failed; returned (%v).",
	15:   "Exit  Initialize(); senzingConfig.SetLogLevel failed; returned (%v).",
	16:   "Exit  Initialize(); senzingConfig.InitializeSenzing; returned (%v).",
	19:   "Exit  Initialize() returned (%v).",
	20:   "Enter InitializeSpecificDatabase().",
	21:   "Exit  InitializeSpecificDatabase(); json.Marshal failed; returned (%v).",
	22:   "Exit  InitializeSpecificDatabase(); engineconfigurationjsonparser.New failed; returned (%v).",
	23:   "Exit  InitializeSpecificDatabase(); parser.GetDatabaseUrls failed; returned (%v).",
	24:   "Exit  InitializeSpecificDatabase(); url.Parse failed; returned (%v).",
	25:   "Exit  InitializeSpecificDatabase(); initializerImpl.initializeSpecificDatabaseSqlite; returned (%v).",
	29:   "Exit  InitializeSpecificDatabase() returned (%v).",
	30:   "Enter RegisterObserver().",
	31:   "Exit  RegisterObserver(); json.Marshal failed; returned (%v).",
	32:   "Exit  RegisterObserver(); initializerImpl.observers.RegisterObserver failed; returned (%v).",
	33:   "Exit  RegisterObserver(); initializerImpl.getSenzingConfig().RegisterObserver failed; returned (%v).",
	34:   "Exit  RegisterObserver(); initializerImpl.getSenzingSchema().RegisterObserver; returned (%v).",
	39:   "Exit  RegisterObserver() returned (%v).",
	40:   "Enter SetLogLevel().",
	41:   "Exit  SetLogLevel(); json.Marshal failed; returned (%v).",
	42:   "Exit  SetLogLevel(); logging.IsValidLogLevelName failed; returned (%v).",
	43:   "Exit  SetLogLevel(); initializerImpl.getLogger().SetLogLevel failed; returned (%v).",
	44:   "Exit  SetLogLevel(); initializerImpl.senzingConfigSingleton.SetLogLevel failed; returned (%v).",
	45:   "Exit  SetLogLevel(); initializerImpl.getSenzingSchema().SetLogLevel failed; returned (%v).",
	49:   "Exit  SetLogLevel() returned (%v).",
	50:   "Enter UnregisterObserver().",
	51:   "Exit  UnregisterObserver(); json.Marshal failed; returned (%v).",
	52:   "Exit  UnregisterObserver(); initializerImpl.getSenzingConfig().UnregisterObserver failed; returned (%v).",
	53:   "Exit  UnregisterObserver(); initializerImpl.getSenzingSchema().UnregisterObserver failed; returned (%v).",
	54:   "Exit  UnregisterObserver(); initializerImpl.observers.UnregisterObserver failed; returned (%v).",
	59:   "Exit  UnregisterObserver() returned (%v).",
	1000: "Initialize parameters: %+v",
	1001: "InitializeSpecificDatabase parameters: %+v",
	1002: "RegisterObserver parameters: %+v",
	1003: "SetLogLevel parameters: %+v",
	1004: "UnregisterObserver parameters: %+v",
	1011: "Initialize(); json.Marshal failed; Error: %v.",
	1012: "Initialize(); initializerImpl.InitializeSpecificDatabase failed; Error: %v.",
	1013: "Initialize(); senzingSchema.SetLogLevel failed; Error: %v.",
	1014: "Initialize(); senzingSchema.InitializeSenzing failed; Error: %v.",
	1015: "Initialize(); senzingConfig.SetLogLevel failed; Error: %v.",
	1016: "Initialize(); senzingConfig.InitializeSenzing; Error: %v.",
	2001: "Created file: %s",
	8001: "Initialize",
	8002: "RegisterObserver",
	8003: "SetLogLevel",
	8004: "UnregisterObserver",
	8005: "InitializeFiles",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
