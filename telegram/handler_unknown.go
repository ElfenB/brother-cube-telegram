package telegram

import (
	"brother-cube-telegram/logger"
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func unknownCommandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	// Check if the message starts with a command (/)
	if !strings.HasPrefix(update.Message.Text, "/") {
		return
	}

	// Extract the command part (without arguments)
	parts := strings.Fields(update.Message.Text)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	logger.Info("Unknown command received: %s from %s", command, update.Message.From.Username)

	// Send error message for unknown commands
	helpText := `❌ Unknown command: ` + command

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpText,
	})
}
