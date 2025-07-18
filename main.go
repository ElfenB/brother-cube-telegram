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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	printer := printers.NewPrinter()

	// Add printer to context (simple approach)
	ctx = context.WithValue(ctx, "printer", printer)

	b := telegram.GetBot()

	log.Println("Bot started successfully. Press Ctrl+C to stop.")
	b.Start(ctx)
}
