package telegram

import (
	"context"
	"log"
	"runtime/debug"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// recoveryMiddleware catches panics and prevents the application from crashing
func recoveryMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic with stack trace
				log.Printf("PANIC recovered: %v\nStack trace:\n%s", r, debug.Stack())

				// Try to send a generic error message to the user
				if update.Message != nil {
					_, err := b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   "❌ An unexpected error occurred. The bot is still running and you can try again.",
					})
					if err != nil {
						log.Printf("Failed to send error message after panic: %v", err)
					}
				}
			}
		}()

		// Call the next handler
		next(ctx, b, update)
	}
}
