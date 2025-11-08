package maxbot

import (
	"context"

	"github.com/Hirogava/CampusShedule/internal/config/logger"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func StartListening(api *maxbot.Api, ctx context.Context) {
	for upd := range api.GetUpdates(ctx) {
		api.Debugs.Send(ctx, upd)

		switch upd := upd.(type) {
			case *schemes.MessageCreatedUpdate:
				out := "bot прочитал текст: " + upd.Message.Body.Text
				switch upd.GetCommand() {
					case "/start":
				}
			default:
				logger.Logger.Println("Unknown message type")
		}
	}
}
