package telegram

import (
	"brother-cube-telegram/utils"
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// printerMiddlewareHandler handles the actual printer checking logic
func printerMiddlewareHandler(mainCtx context.Context, next bot.HandlerFunc) bot.HandlerFunc {
	return func(handlerCtx context.Context, b *bot.Bot, update *models.Update) {
		// Get printer from the main context
		printer := utils.GetPrinterFromContext(mainCtx)

		if printer != nil {
			// Printer is available - add it to handler context and proceed
			handlerCtx = context.WithValue(handlerCtx, "printer", printer)
			next(handlerCtx, b, update)
			return
		}

		// Printer not available - send error message
		if update.Message != nil {
			_, err := b.SendMessage(handlerCtx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "❌ Printer is not available. Please check the printer connection.",
			})
			if err != nil {
				log.Printf("Error sending printer unavailable message: %v", err)
			}
		}
		// Don't proceed to handlers when printer is not available
	}
}

// Create a middleware factory that captures the main context
func createMiddlewareWithCtxFactory(mainCtx context.Context, middlewareHandler func(context.Context, bot.HandlerFunc) bot.HandlerFunc) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return middlewareHandler(mainCtx, next)
	}
}
