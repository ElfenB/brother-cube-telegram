package telegram

import (
	"context"
	"log"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"brother-cube-telegram/utils"
)

func GetBot() *bot.Bot {
	// Get bot token from environment variable
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	return b
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Println("Received message:", update.Message.Text, "from", update.Message.From.Username)

	// Get printer from context
	printer := utils.GetPrinterFromContext(ctx)

	if printer == nil {
		log.Println("No printer found in context")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error: The printer is not available.",
		})
		return
	}

	version := printer.GetVersion()
	log.Printf("Printer version: %s", version)

	info, err := printer.GetPrinterInfo()
	if err != nil {
		log.Printf("Error getting printer info: %v", err)
	} else {
		log.Printf("Printer info: %s", info)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
