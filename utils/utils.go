package utils

import (
	"brother-cube-telegram/printers"
	"context"
	"log"
)

const printerCtxKey string = "printer"

// GetPrinterFromContext gets printer from context
func GetPrinterFromContext(ctx context.Context) *printers.Printer {
	if printer, ok := ctx.Value(printerCtxKey).(*printers.Printer); ok {
		return printer
	}
	log.Println("Warning: Printer not found in context")
	return nil
}
