package senzingload

import (
	"context"
	"errors"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type SenzingLoad interface {
	LoadURLs(ctx context.Context) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identifier of the  package found messages having the format "senzing-6504xxxx".
const ComponentID = 6504

// Log message prefix.
const Prefix = "init-database.senzingload."

const (
	OptionCallerSkip4 = 4
	OptionCallerSkip5 = 5
)

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for senzingload implementation.
var IDMessages = map[int]string{
	10:   "Enter " + Prefix + "LoadURLs().",
	29:   "Exit  " + Prefix + "LoadURLs(); copyFile when replacing template/szConfig.json failed; returned (%v).",
	30:   "Enter " + Prefix + "RegisterObserver(%s).",
	31:   "Exit  " + Prefix + "RegisterObserver(%s); json.Marshal failed; returned (%v).",
	32:   "Exit  " + Prefix + "RegisterObserver(%s); senzingLoad.observers.RegisterObserver failed; returned (%v).",
	33:   "Exit  " + Prefix + "RegisterObserver(%s); senzingLoad.getDependentServices failed; returned (%v).",
	34:   "Exit  " + Prefix + "RegisterObserver(%s); szConfig.RegisterObserver failed; returned (%v).",
	35:   "Exit  " + Prefix + "RegisterObserver(%s); szConfigmgr.RegisterObserver failed; returned (%v).",
	39:   "Exit  " + Prefix + "RegisterObserver(%s) returned (%v).",
	40:   "Enter " + Prefix + "SetLogLevel(%s).",
	41:   "Exit  " + Prefix + "SetLogLevel(%s); json.Marshal failed; returned (%v).",
	42:   "Exit  " + Prefix + "SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).",
	43:   "Exit  " + Prefix + "SetLogLevel(%s); senzingLoad.getLogger().SetLogLevel failed; returned (%v).",
	44:   "Exit  " + Prefix + "SetLogLevel(%s); senzingLoad.getDependentServices failed; returned (%v).",
	45:   "Exit  " + Prefix + "SetLogLevel(%s); szConfig.SetLogLevel failed; returned (%v).",
	46:   "Exit  " + Prefix + "SetLogLevel(%s); szConfigmgr.SetLogLevel failed; returned (%v).",
	49:   "Exit  " + Prefix + "SetLogLevel(%s) returned (%v).",
	50:   "Enter " + Prefix + "UnregisterObserver(%s).",
	51:   "Exit  " + Prefix + "UnregisterObserver(%s); json.Marshal failed; returned (%v).",
	52:   "Exit  " + Prefix + "UnregisterObserver(%s); szConfig.UnregisterObserver failed; returned (%v).",
	53:   "Exit  " + Prefix + "UnregisterObserver(%s); szConfigmgr.UnregisterObserver failed; returned (%v).",
	54:   "Exit  " + Prefix + "UnregisterObserver(%s); senzingLoad.observers.UnregisterObserver failed; returned (%v).",
	59:   "Exit  " + Prefix + "UnregisterObserver(%s) returned (%v).",
	60:   "Enter " + Prefix + "SetObserverOrigin(%s).",
	61:   "Exit  " + Prefix + "SetObserverOrigin(%s); json.Marshal failed; returned (%v).",
	69:   "Exit  " + Prefix + "SetObserverOrigin(%s).",
	1001: Prefix + "LoadURLs parameters: %+v",
	1002: Prefix + "RegisterObserver parameters: %+v",
	1003: Prefix + "SetLogLevel parameters: %+v",
	1004: Prefix + "SetObserverOrigin parameters: %+v",
	1005: Prefix + "UnregisterObserver parameters: %+v",
	2002: "Processed %d records for URL: %s",
	2003: "Processed %d redo records",
	3001: "Processing URL: %s",
	8001: Prefix + "LoadURLs",
	8003: Prefix + "RegisterObserver",
	8004: Prefix + "SetLogLevel",
	8005: Prefix + "SetObserverOrigin",
	8006: Prefix + "UnregisterObserver",
}

// Status strings for specific messages.
var IDStatuses = map[int]string{}

var errForPackage = errors.New("senzingload")
