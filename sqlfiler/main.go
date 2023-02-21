package sqlfiler

import (
	"context"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-observing/observer"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type SqlFiler interface {
	ProcessFile(ctx context.Context) error
	RegisterObserver(ctx context.Context, observer observer.Observer) error
	SetLogLevel(ctx context.Context, logLevel logger.Level) error
	UnregisterObserver(ctx context.Context, observer observer.Observer) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-6999xxxx".
const ProductId = 6999

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for sqlfiler implementation.
var IdMessages = map[int]string{
	1:    "Enter InitializeSenzingConfiguration().",
	2:    "Exit  InitializeSenzingConfiguration() returned (%v).",
	2000: "Entry: %+v",
	4001: "Call to net.Listen(tcp, %s) failed.",
	5001: "Failed to serve.",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
