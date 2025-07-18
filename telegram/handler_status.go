package telegram

import (
	"brother-cube-telegram/utils"
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func statusHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	printer := utils.GetPrinterFromContext(ctx)

	// Use the printer to get the status
	status, err := printer.GetPrinterInfo()
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error getting printer status: " + err.Error(),
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Printer status: " + status,
	})
}
