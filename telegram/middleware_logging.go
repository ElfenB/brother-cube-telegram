package telegram

import (
	"brother-cube-telegram/logger"
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func loggingMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			logger.Debug("Received message: %s from %s", update.Message.Text, update.Message.From.Username)
		}
		next(ctx, b, update)
	}
}
