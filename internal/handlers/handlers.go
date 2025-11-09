package handlers

import (
	"context"
	"strconv"

	"github.com/Hirogava/CampusShedule/internal/config/logger"
	"github.com/Hirogava/CampusShedule/internal/repository/postgres"
	keyboard "github.com/Hirogava/CampusShedule/internal/service/maxbot"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
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

func ScheduleHandler(api *maxbot.Api, upd *schemes.MessageCallbackUpdate, manager *postgres.Manager, ctx context.Context) {
	userID := upd.Message.Sender.UserId

	hasUniversity, err := manager.HasUserUniversity(userID)
	if err != nil {
		logger.Logger.Printf("Error checking user university: %v", err)
		return
	}

	if hasUniversity {
		// жестко присылаю расписание
		return
	}

	out := "Выберите университет"
	kb := keyboard.CreateKeyboardForUniversities(api, manager)
	helpingErrorCallbackSending(api, upd, kb, out, ctx)
}

func HandleUniversitySelection(api *maxbot.Api, upd *schemes.MessageCallbackUpdate, manager *postgres.Manager, ctx context.Context, unID string) {
	intID, _ := strconv.Atoi(unID)

	_, err := manager.SetUserUniversity(upd.Callback.User.UserId, intID)
	if err != nil {
		logger.Logger.Printf("Error setting user university: %v", err)
		return
	}

	// TODO:
	// тут надо будет сделать функцию запроса групп университета (апи дается первым в прошлой функции)
	var groups []string // пока просто заглушка

	out := "Выберите группу"
	kb := keyboard.CreateKeyboardForGroups(api, groups)
	helpingErrorCallbackSending(api, upd, kb, out, ctx)
}

func HandleGroupSelection(api *maxbot.Api, upd *schemes.MessageCallbackUpdate, manager *postgres.Manager, ctx context.Context, group string) {
	err := manager.SetUserGroup(upd.Callback.User.UserId, group)
	if err != nil {
		logger.Logger.Printf("Error setting user group: %v", err)
		return
	}

	out := "Поздравляю блаблабла..."
	kb := keyboard.CreateKeyboardForStart(api)
	helpingErrorCallbackSending(api, upd, kb, out, ctx)
}

func helpingErrorCallbackSending(api *maxbot.Api, upd *schemes.MessageCallbackUpdate, kb *maxbot.Keyboard, out string, ctx context.Context) {
	res, err := api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).AddKeyboard(kb).SetText(out))
	if err != nil {
		logger.Logger.Printf("Error sending message: %v", err)
	} else {
		logger.Logger.Printf("Message sent: %v", res)
	}
}
