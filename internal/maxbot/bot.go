package maxbot

import (
	"context"

	"github.com/Hirogava/CampusShedule/internal/config/logger"
	"github.com/Hirogava/CampusShedule/internal/handlers"
	"github.com/Hirogava/CampusShedule/internal/models/buttons"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func StartListening(api *maxbot.Api, ctx context.Context) {
	for upd := range api.GetUpdates(ctx) {
		api.Debugs.Send(ctx, upd)

		switch upd := upd.(type) {
			case *schemes.MessageCreatedUpdate:
				out := "bot прочитал текст: " + upd.Message.Body.Text // текст сообщения
				switch upd.GetCommand() {
					case "/start":
						handlers.StartHandler(api, upd, ctx)
						continue
				}
			case *schemes.MessageCallbackUpdate:
				msg := maxbot.NewMessage()
				if upd.Message.Recipient.UserId != 0 {
					msg.SetUser(upd.Message.Recipient.UserId)
				}
				switch upd.Callback.Payload {
				case string(buttons.BtnSchedule):
				case string(buttons.BtnDecan):
				case string(buttons.BtnProjects):
				}
			default:
				logger.Logger.Println("Unknown message type")
		}
	}
}
