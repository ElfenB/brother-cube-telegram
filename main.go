package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"

	"brother-cube-telegram/gpio"
	"brother-cube-telegram/printers"
	"brother-cube-telegram/telegram"
)

func main() {
	// Add recovery for the main function
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Main function panic recovered: %v", r)
		}
	}()

	relay, err := gpio.NewRelay(17)
	if err != nil {
		log.Printf("Failed to initialize relay: %v", err)
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

	log.Println("Bot started successfully. Press Ctrl+C to stop.")

	// Start the bot (this blocks until context is cancelled)
	b.Start(ctx)

	log.Println("Bot stopped gracefully.")
}
