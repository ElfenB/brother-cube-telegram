package telegram

import (
	"brother-cube-telegram/logger"
	"context"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func authorizationMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		// Skip authorization if no message is present
		if update.Message == nil {
			next(ctx, b, update)
			return
		}

		chatID := update.Message.Chat.ID

		// Get allowed chat IDs from environment variable
		allowedChatIDsStr := os.Getenv("TELEGRAM_ALLOWED_CHAT_IDS")
		if allowedChatIDsStr == "" {
			logger.Debug("TELEGRAM_ALLOWED_CHAT_IDS environment variable not set - allowing all access")
			next(ctx, b, update)
			return
		}

		// Parse the comma-separated list of allowed chat IDs
		allowedChatIDs := parseAllowedChatIDs(allowedChatIDsStr)

		// Check if the current chat ID is in the allowed list
		if !isAuthorized(chatID, allowedChatIDs) {
			logger.Warn("Unauthorized access attempt from chat ID: %d", chatID)
			sendUnauthorizedMessage(ctx, b, chatID)
			return
		}

		logger.Debug("Authorized access from chat ID: %d", chatID)
		next(ctx, b, update)
	}
}

// Parses a comma-separated string of chat IDs into a slice of int64
func parseAllowedChatIDs(chatIDsStr string) []int64 {
	var allowedChatIDs []int64

	// Split by comma and parse each ID
	chatIDStrings := strings.SplitSeq(chatIDsStr, ",")
	for idStr := range chatIDStrings {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		chatID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			logger.Error("Invalid chat ID in TELEGRAM_ALLOWED_CHAT_IDS: %s", idStr)
			continue
		}

		allowedChatIDs = append(allowedChatIDs, chatID)
	}

	return allowedChatIDs
}

// Checks if a chat ID is in the list of allowed chat IDs
func isAuthorized(chatID int64, allowedChatIDs []int64) bool {
	return slices.Contains(allowedChatIDs, chatID)
}

// Sends an unauthorized access message to the user
func sendUnauthorizedMessage(ctx context.Context, b *bot.Bot, chatID int64) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "🚫 Unauthorized access. This bot is restricted to specific users only.",
	})
	if err != nil {
		logger.Error("Failed to send unauthorized message to chat ID %d: %v", chatID, err)
	}
}
