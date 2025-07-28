package main

import (
	"context"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"

	"brother-cube-telegram/gpio"
	"brother-cube-telegram/logger"
	"brother-cube-telegram/printers"
	"brother-cube-telegram/telegram"
)

func main() {
	// Enable debug logging to see command executions
	logger.SetLogLevel(logger.DEBUG)

	// Add recovery for the main function
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Main function panic recovered: %v", r)
		}
	}()

	logger.Info("Testing relay on GPIO17...")
	relay, err := gpio.NewRelay(17)
	if err != nil {
		logger.Error("Failed to initialize relay: %v", err)
	} else {
		defer relay.Close()
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	printer := printers.NewPrinter(relay)
	if printer != nil {
		defer printer.Close()
	}

	// Add printer to context
	ctx = context.WithValue(ctx, "printer", printer)

	b := telegram.GetBot(ctx)

	logger.Info("Bot started successfully. Press Ctrl+C to stop.")

	// Start the bot (this blocks until context is cancelled)
	b.Start(ctx)

	logger.Info("Bot stopped gracefully.")
}
