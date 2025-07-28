package main

import (
	"context"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"

	"brother-cube-telegram/config"
	"brother-cube-telegram/gpio"
	"brother-cube-telegram/logger"
	"brother-cube-telegram/printers"
	"brother-cube-telegram/telegram"
)

func main() {
	// Load configuration
	if err := config.LoadDefault(); err != nil {
		logger.Error("Failed to load configuration: %v", err)
		return
	}
	logger.Info("Configuration loaded successfully")

	// Set log level from configuration
	cfg := config.Get()
	logger.SetLogLevel(cfg.Logging.GetLogLevel())

	// Add recovery for the main function
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Main function panic recovered: %v", r)
		}
	}()

	logger.Info("Testing relay on GPIO%d...", cfg.GPIO.RelayPin)
	relay, err := gpio.NewRelay(cfg.GPIO.RelayPin)
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
