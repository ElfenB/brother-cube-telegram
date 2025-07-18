package telegram

import (
	"brother-cube-telegram/utils"
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Defer a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in defaultHandler: %v", r)

			// Try to send an error message to the user if possible
			if update.Message != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "❌ An error occurred while processing your message. Please try again.",
				})
			}
		}
	}()

	// Check if update has a message and message has text
	if update.Message == nil {
		log.Println("Received update without message")
		return
	}

	if update.Message.Text == "" {
		log.Println("Received message without text (possibly media)")
		return
	}

	if update.Message.From == nil {
		log.Println("Received message without sender info")
		return
	}

	log.Println("Received message:", update.Message.Text, "from", update.Message.From.Username)

	printer := utils.GetPrinterFromContext(ctx)

	// Wrap the printer call in error handling
	err := printer.PrintLabelYolo(update.Message.Text)
	if err != nil {
		log.Printf("Error printing label: %v", err)

		// Inform the user about the error
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Failed to print label: " + err.Error(),
		})
		return
	}

	// Send success message
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "✅ Label printed successfully!",
	})
}
