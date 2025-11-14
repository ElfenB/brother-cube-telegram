package telegram

import (
	"brother-cube-telegram/config"
	"brother-cube-telegram/logger"
	"brother-cube-telegram/utils"
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func ppreviewHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Defer a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic in ppreviewHandler: %v", r)

			// Try to send an error message to the user if possible
			if update.Message != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "❌ An error occurred while processing your preset preview command. Please try again.",
				})
			}
		}
	}()

	// Check if update has a message and message has text
	if update.Message == nil {
		logger.Warn("Received update without message in ppreviewHandler")
		return
	}

	if update.Message.Text == "" {
		logger.Debug("Received message without text in ppreviewHandler")
		return
	}

	if update.Message.From == nil {
		logger.Warn("Received message without sender info in ppreviewHandler")
		return
	}

	// Parse the command: /ppreview <preset_name> <text_to_preview>
	parts := strings.SplitN(strings.TrimSpace(update.Message.Text), " ", 3)

	if len(parts) < 2 {
		// Show available presets (reuse the same function as preset handler)
		sendPresetUsage(ctx, b, update.Message.Chat.ID)
		return
	}

	if len(parts) < 3 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   GetCommandUsageMessageWithError("ppreview", "Missing text to preview."),
		})
		return
	}

	presetName := parts[1]
	textToPreview := parts[2]

	// Get the configuration and printer
	cfg := config.Get()
	printer := utils.GetPrinterFromContext(ctx)

	// Look up the preset
	preset := cfg.Printer.GetPreset(presetName)
	if preset == nil {
		sendPresetNotFound(ctx, b, update.Message.Chat.ID, presetName)
		return
	}

	logger.Info("Using preset '%s' (font size: %d, font family: %s) to preview: %s from %s",
		presetName, preset.FontSize, preset.FontFamily, textToPreview, update.Message.From.Username)

	// Generate preview with the preset's font size and family
	img, err := printer.PreviewLabelWithPreset(textToPreview, update.Message.From.ID, preset)
	if err != nil {
		logger.Error("Error generating preset preview for '%s': %v", presetName, err)

		// Inform the user about the error
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ Failed to generate preview with preset '%s': %s", presetName, err.Error()),
		})
		return
	}

	// Send the preview image
	fontInfo := fmt.Sprintf("📏 Font size: %d", preset.FontSize)
	if preset.FontFamily != "" {
		fontInfo += fmt.Sprintf(", Font: %s", preset.FontFamily)
	}

	caption := fmt.Sprintf("Preview using preset '%s'\n%s", presetName, fontInfo)

	b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:  update.Message.Chat.ID,
		Photo:   &models.InputFileUpload{Filename: "preset-preview.png", Data: bytes.NewReader(img)},
		Caption: caption,
	})
}
