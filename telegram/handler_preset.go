package telegram

import (
	"brother-cube-telegram/config"
	"brother-cube-telegram/logger"
	"brother-cube-telegram/utils"
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func presetHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Defer a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic in presetHandler: %v", r)

			// Try to send an error message to the user if possible
			if update.Message != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "❌ An error occurred while processing your preset command. Please try again.",
				})
			}
		}
	}()

	// Check if update has a message and message has text
	if update.Message == nil {
		logger.Warn("Received update without message in presetHandler")
		return
	}

	if update.Message.Text == "" {
		logger.Debug("Received message without text in presetHandler")
		return
	}

	if update.Message.From == nil {
		logger.Warn("Received message without sender info in presetHandler")
		return
	}

	// Parse the command: /preset <preset_name> <text_to_print>
	parts := strings.SplitN(strings.TrimSpace(update.Message.Text), " ", 3)

	if len(parts) < 2 {
		// Show available presets
		sendPresetUsage(ctx, b, update.Message.Chat.ID)
		return
	}

	if len(parts) < 3 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Usage: /preset <preset_name> <text_to_print>\n\nExample: /preset kitchen my text",
		})
		return
	}

	presetName := parts[1]
	textToPrint := parts[2]

	// Get the configuration and printer
	cfg := config.Get()
	printer := utils.GetPrinterFromContext(ctx)

	// Look up the preset
	preset := cfg.Printer.GetPreset(presetName)
	if preset == nil {
		sendPresetNotFound(ctx, b, update.Message.Chat.ID, presetName)
		return
	}

	logger.Info("Using preset '%s' (font size: %d, font family: %s) to print: %s from %s",
		presetName, preset.FontSize, preset.FontFamily, textToPrint, update.Message.From.Username)

	// Print the label with the preset's font size and family
	err := printer.PrintLabelWithPreset(textToPrint, preset)
	if err != nil {
		logger.Error("Error printing label with preset '%s': %v", presetName, err)

		// Inform the user about the error
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ Failed to print label with preset '%s': %s", presetName, err.Error()),
		})
		return
	}

	// Send success message
	fontInfo := fmt.Sprintf("📏 Font size: %d", preset.FontSize)
	if preset.FontFamily != "" {
		fontInfo += fmt.Sprintf(", Font: %s", preset.FontFamily)
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("✅ Label printed successfully using preset '%s'!\n%s", presetName, fontInfo),
	})
}

func sendPresetUsage(ctx context.Context, b *bot.Bot, chatID int64) {
	cfg := config.Get()
	presetNames := cfg.Printer.GetPresetNames()

	if len(presetNames) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ No presets are configured.",
		})
		return
	}

	var message strings.Builder
	message.WriteString("📋 Available presets:\n\n")

	for _, name := range presetNames {
		preset := cfg.Printer.GetPreset(name)
		if preset != nil {
			fontInfo := fmt.Sprintf("font size: %d", preset.FontSize)
			if preset.FontFamily != "" {
				fontInfo += fmt.Sprintf(", font: %s", preset.FontFamily)
			}
			message.WriteString(fmt.Sprintf("• **%s** - %s (%s)\n", name, preset.Description, fontInfo))
		}
	}

	message.WriteString("\n💡 Usage: `/preset <preset_name> <text_to_print>`\n")
	message.WriteString("Example: `/preset kitchen my text`")

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      message.String(),
		ParseMode: models.ParseModeMarkdown,
	})
}

func sendPresetNotFound(ctx context.Context, b *bot.Bot, chatID int64, presetName string) {
	cfg := config.Get()
	presetNames := cfg.Printer.GetPresetNames()

	var message strings.Builder
	message.WriteString(fmt.Sprintf("❌ Preset '%s' not found.\n\n", presetName))

	if len(presetNames) > 0 {
		message.WriteString("📋 Available presets:\n")
		for _, name := range presetNames {
			preset := cfg.Printer.GetPreset(name)
			if preset != nil {
				message.WriteString(fmt.Sprintf("• **%s** - %s\n", name, preset.Description))
			}
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      message.String(),
		ParseMode: models.ParseModeMarkdown,
	})
}
