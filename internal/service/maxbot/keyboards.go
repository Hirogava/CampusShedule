package maxbot

import (
	"github.com/Hirogava/CampusShedule/internal/models/buttons"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func CreateKeyboardForStart(api *maxbot.Api) *maxbot.Keyboard {
	keyboard := api.Messages.NewKeyboardBuilder()
	keyboard.
		AddRow().
		AddCallback("Узнать расписание", schemes.POSITIVE, string(buttons.BtnSchedule))
	keyboard.
		AddRow().
		AddCallback("Электронный деканат", schemes.POSITIVE, string(buttons.BtnDecan))
	keyboard.
		AddRow().
		AddCallback("Проекты твоего вуза", schemes.POSITIVE, string(buttons.BtnProjects))

	return keyboard
}
