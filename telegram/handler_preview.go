package telegram

import (
	"brother-cube-telegram/logger"
	"brother-cube-telegram/utils"
	"bytes"
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func previewHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Safety checks
	if update.Message == nil || update.Message.Text == "" {
		logger.Warn("Preview handler: Invalid message")
		return
	}

	printer := utils.GetPrinterFromContext(ctx)

	// Parse command arguments
	parts := strings.SplitN(strings.TrimSpace(update.Message.Text), " ", 2)

	if len(parts) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   GetCommandUsageMessage("preview"),
		})
		return
	}

	rawText := strings.TrimSpace(parts[1])

	if rawText == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   GetCommandUsageMessageWithError("preview", "Text cannot be empty."),
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
