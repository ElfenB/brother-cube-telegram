package telegram

import (
	"brother-cube-telegram/logger"
	"context"
	"os"

	"github.com/go-telegram/bot"
)

func GetBot(ctx context.Context) *bot.Bot {
	// Get bot token from environment variable
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		logger.Error("TELEGRAM_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMiddlewares(
			recoveryMiddleware,
			createMiddlewareWithCtxFactory(ctx, printerMiddlewareHandler),
			loggingMiddleware,
		),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		logger.Error("Failed to create bot: %v", err)
		os.Exit(1)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "status", bot.MatchTypeCommandStartOnly, statusHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "preview", bot.MatchTypeCommand, previewHandler)

	return b
}
