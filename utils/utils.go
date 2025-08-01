package utils

import (
	"brother-cube-telegram/gpio"
	"brother-cube-telegram/logger"
	"brother-cube-telegram/printers"
	"context"
)

const printerCtxKey string = "printer"
const relayCtxKey string = "relay"

// GetPrinterFromContext gets printer from context
func GetPrinterFromContext(ctx context.Context) *printers.Printer {
	if printer, ok := ctx.Value(printerCtxKey).(*printers.Printer); ok {
		return printer
	}
	logger.Warn("Printer not found in context")
	return nil
}

// GetRelayFromContext gets relay from context
func GetRelayFromContext(ctx context.Context) *gpio.Relay {
	if relay, ok := ctx.Value(relayCtxKey).(*gpio.Relay); ok {
		return relay
	}
	logger.Warn("Relay not found in context")
	return nil
}
