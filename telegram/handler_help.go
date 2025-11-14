package telegram

import (
	"brother-cube-telegram/logger"
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Defer a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic in helpHandler: %v", r)

			// Try to send an error message to the user if possible
			if update.Message != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "❌ An error occurred while processing your help request. Please try again.",
				})
			}
		}
	}()

	// Check if update has a message and message has text
	if update.Message == nil {
		logger.Warn("Received update without message in helpHandler")
		return
	}

	if update.Message.From == nil {
		logger.Warn("Received message without sender info in helpHandler")
		return
	}

	// Parse the command: /help [command_name]
	parts := strings.SplitN(strings.TrimSpace(update.Message.Text), " ", 2)

	logger.Info("Help requested by %s", update.Message.From.Username)

	// Check if specific command help was requested
	if len(parts) > 1 {
		commandName := strings.TrimSpace(parts[1])
		// Remove leading slash unconditionally
		commandName = strings.TrimPrefix(commandName, "/")

		logger.Debug("Specific help requested for command: %s", commandName)

		// Send specific command help
		helpText := GetCommandUsageMessage(commandName)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   helpText,
		})

		if err != nil {
			logger.Error("Failed to send specific command help: %v", err)
		}
		return
	}

	// Build the general help message
	var message strings.Builder
	message.WriteString("🤖 Brother Cube Telegram Bot Help\n\n")
	message.WriteString("Available commands:\n\n")

	commands := GetRegisteredCommands()

	for _, cmd := range commands {
		message.WriteString(fmt.Sprintf("%s\n", cmd.Command))
		message.WriteString(fmt.Sprintf("• %s\n", cmd.Description))
		message.WriteString(fmt.Sprintf("• Usage: %s\n", cmd.Usage))
		message.WriteString(fmt.Sprintf("• Example: %s\n\n", cmd.Example))
	}

	message.WriteString("📝 Tips:\n")
	message.WriteString("• Use /preset without arguments to see available presets\n")
	message.WriteString("• Preview your labels before printing to save tape\n")
	message.WriteString("• The bot will automatically manage printer power\n")
	message.WriteString("• Use '/help <command>' for detailed help on a specific command\n\n")

	// Send the help message
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   message.String(),
	})

	if err != nil {
		logger.Error("Failed to send help message: %v", err)
		// Try sending a simple fallback message
		_, fallbackErr := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Help message failed to send. Available commands: /help, /status, /preview, /size, /preset\n\nFor detailed help, contact support.",
		})
		if fallbackErr != nil {
			logger.Error("Failed to send fallback help message: %v", fallbackErr)
		}
	}
}
