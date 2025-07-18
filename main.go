package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"

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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	printer := printers.NewPrinter()

	// Add printer to context (simple approach)
	ctx = context.WithValue(ctx, "printer", printer)

	b := telegram.GetBot(ctx)

	log.Println("Bot started successfully. Press Ctrl+C to stop.")

	// Start the bot (this blocks until context is cancelled)
	b.Start(ctx)

	log.Println("Bot stopped gracefully.")
}
