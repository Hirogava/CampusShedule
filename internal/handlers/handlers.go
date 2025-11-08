package handlers

import (
	"context"

	"github.com/Hirogava/CampusShedule/internal/config/logger"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	keyboard "github.com/Hirogava/CampusShedule/internal/service/maxbot"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func StartHandler(api *maxbot.Api, upd *schemes.MessageCreatedUpdate, ctx context.Context) {
	out := "Приветствую блаблаб"
	kb := keyboard.CreateKeyboardForStart(api)
	res, err := api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).AddKeyboard(kb).SetText(out))
	if err != nil {
		logger.Logger.Printf("Error sending message: %v", err)
	} else {
		logger.Logger.Printf("Message sent: %v", res)
	}
}
