package maxbot

import (
	"context"
	"strings"

	"github.com/Hirogava/CampusShedule/internal/config/logger"
	"github.com/Hirogava/CampusShedule/internal/handlers"
	"github.com/Hirogava/CampusShedule/internal/models/buttons"
	"github.com/Hirogava/CampusShedule/internal/repository/postgres"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func StartListening(api *maxbot.Api, manager *postgres.Manager, ctx context.Context) {
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

				go HandleCallback(api, upd, manager, ctx)

				switch upd.Callback.Payload {
				case string(buttons.BtnSchedule):
					handlers.ScheduleHandler(api, upd, manager, ctx)
				case string(buttons.BtnDecan):
				case string(buttons.BtnProjects):
				}
			default:
				logger.Logger.Println("Unknown message type: ", upd)
		}
	}
}

func HandleCallback(api *maxbot.Api, upd *schemes.MessageCallbackUpdate, manager *postgres.Manager, ctx context.Context) {
	data := upd.Callback.Payload

	if strings.HasPrefix(data, "uni:") {
		university := strings.TrimPrefix(data, "uni:")
		handlers.HandleUniversitySelection(api, upd, manager, ctx, university)
		return
	}

	if strings.HasPrefix(data, "group:") {
		group := strings.TrimPrefix(data, "group:")
		handlers.HandleGroupSelection(api, upd, manager, ctx, group)
		return
	}
}
