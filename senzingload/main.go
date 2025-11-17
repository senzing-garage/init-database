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

// Identfier of the  package found messages having the format "senzing-6502xxxx".
const ComponentID = 6504

// Log message prefix.
const Prefix = "init-database.senzingconfig."

const (
	OptionCallerSkip4 = 4
	OptionCallerSkip5 = 5
)

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for sqlfiler implementation.
var IDMessages = map[int]string{
	10: "Enter " + Prefix + "LoadURLs().",
	24: "Exit  " + Prefix + "InitializeSenzing(); copyFile when replacing template/szConfig.json failed; returned (%v).",
	29: "Exit  " + Prefix + "InitializeSenzing() returned (%v).",

	1001: Prefix + "InitializeSenzing parameters: %+v",
	1002: Prefix + "RegisterObserver parameters: %+v",

	2001: "Processing URL: %s",
	2002: "Processed %d records for URL: %s",
	2003: "Processed %d redo records",
	4001: "When comparing %s and %s, an error occurred. Assuming files not equal.",
	5001: "File does not exist: %s [SENZING_TOOLS_ENGINE_CONFIGURATION_FILE]",
	5002: "Could not backup %s to %s",
	5003: "Could not copy %s to %s [SENZING_TOOLS_ENGINE_CONFIGURATION_FILE]",
	8001: Prefix + "LoadURLs",
}

// Status strings for specific messages.
var IDStatuses = map[int]string{}

var errForPackage = errors.New("senzingload")
