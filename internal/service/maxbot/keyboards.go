package maxbot

import (
	"fmt"

	"github.com/Hirogava/CampusShedule/internal/config/logger"
	"github.com/Hirogava/CampusShedule/internal/models/buttons"
	"github.com/Hirogava/CampusShedule/internal/models/db"
	"github.com/Hirogava/CampusShedule/internal/repository/postgres"
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

func CreateKeyboardForUniversities(api *maxbot.Api, manager *postgres.Manager) *maxbot.Keyboard {
	universities, err := manager.GetUniversities()
	if err != nil {
		logger.Logger.Printf("Error getting universities: %v", err)
		return nil
	}

	kb := api.Messages.NewKeyboardBuilder()
	for _, university := range universities {
		kb.AddRow().AddCallback(university.Name, schemes.POSITIVE, fmt.Sprintf("uni:%d", university.ID))
	}

	return kb
}

func CreateKeyboardForGroups(api *maxbot.Api, groups []db.Group) *maxbot.Keyboard {
	kb := api.Messages.NewKeyboardBuilder()

	for _, group := range groups {
		kb.
			AddRow().
			AddCallback(group.Name, schemes.POSITIVE, fmt.Sprintf("group:%d", group.ID))
	}

	return kb
}
