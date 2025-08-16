package telegram

import (
	"brother-cube-telegram/logger"
	"brother-cube-telegram/utils"
	"context"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func sizeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Defer a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic in sizeHandler: %v", r)

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
		logger.Warn("Received update without message")
		return
	}

	if update.Message.Text == "" {
		logger.Debug("Received message without text (possibly media)")
		return
	}

	if update.Message.From == nil {
		logger.Warn("Received message without sender info")
		return
	}

	logger.Info("Received message: %s from %s", update.Message.Text, update.Message.From.Username)

	printer := utils.GetPrinterFromContext(ctx)

	fontSize := strings.Split(update.Message.Text, " ")[1]
	label := strings.TrimSpace(strings.Join(strings.Split(update.Message.Text, " ")[2:], " "))

	// Try to convert fontSize to int
	fontSizeInt, err := strconv.Atoi(fontSize)
	if err != nil {
		logger.Error("Invalid font size: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Invalid font size. Please provide a valid number.",
		})
		return
	}

	if label == "" {
		logger.Error("Label is empty")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Label is empty. Please provide a valid label after the size information.",
		})
		return
	}

	// Wrap the printer call in error handling
	err = printer.PrintLabel(label, fontSizeInt)
	if err != nil {
		logger.Error("Error printing label: %v", err)

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
