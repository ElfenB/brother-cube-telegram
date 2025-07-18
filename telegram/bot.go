package telegram

import (
	"context"
	"log"
	"os"

	"github.com/go-telegram/bot"
)

func GetBot(ctx context.Context) *bot.Bot {
	// Get bot token from environment variable
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
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
		log.Fatal("Failed to create bot:", err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "status", bot.MatchTypeCommandStartOnly, statusHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "preview", bot.MatchTypeCommand, previewHandler)

	return b
}
