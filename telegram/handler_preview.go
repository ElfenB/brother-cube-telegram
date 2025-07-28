package telegram

import (
	"brother-cube-telegram/logger"
	"brother-cube-telegram/utils"
	"bytes"
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func previewHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Safety checks
	if update.Message == nil || update.Message.Text == "" {
		logger.Warn("Preview handler: Invalid message")
		return
	}

	// Check if message is long enough to contain "preview" command
	if len(update.Message.Text) < 8 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Usage: /preview <text to preview>",
		})
		return
	}

	printer := utils.GetPrinterFromContext(ctx)

	// Remove the "/preview " command from the message text (9 characters)
	rawText := update.Message.Text[9:]

	if rawText == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Usage: /preview <text to preview>",
		})
		return
	}

	img, err := printer.PreviewLabel(rawText, update.Message.From.ID)

	if err != nil {
		logger.Error("Error generating label preview: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error generating label preview: " + err.Error(),
		})
		return
	}

	b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:  update.Message.Chat.ID,
		Photo:   &models.InputFileUpload{Filename: "image.png", Data: bytes.NewReader(img)},
		Caption: "Preview of your label",
	})
}
