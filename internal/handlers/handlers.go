package handlers

import (
	"context"
	"strconv"

	"github.com/Hirogava/CampusShedule/internal/config/logger"
	dbErrors "github.com/Hirogava/CampusShedule/internal/errors/db"
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
		groupID, err := manager.GetUserGroup(userID)
		if err != nil {
			if err == dbErrors.ErrUserNotFound {
				out := "Выберите университет"
				kb := keyboard.CreateKeyboardForUniversities(api, manager)
				helpingErrorCallbackSending(api, upd, kb, out, ctx)
				return
			}
			logger.Logger.Printf("Error getting user group: %v", err)
			return
		}

		schedule, err := manager.GetWeekSchedule(ctx, groupID)
		if err != nil {
			logger.Logger.Printf("Error getting user schedule: %v", err)
			if err == dbErrors.ErrScheduleNotFound {
				out := "У вас нет расписания"
				kb := keyboard.CreateKeyboardForStart(api)
				helpingErrorCallbackSending(api, upd, kb, out, ctx)
				return
			}

			out := "Произошла ошибка"
			kb := keyboard.CreateKeyboardForStart(api)
			helpingErrorCallbackSending(api, upd, kb, out, ctx)
			return
		}

		out := keyboard.CreateScheduledMessage(schedule)
		kb := keyboard.CreateKeyboardForStart(api)
		helpingErrorCallbackSending(api, upd, kb, out, ctx)
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

	groups, err := manager.GetUniversityGroups(intID)
	if err != nil {
		logger.Logger.Printf("Error getting groups: %v", err)
		return
	}

	out := "Выберите группу"
	kb := keyboard.CreateKeyboardForGroups(api, groups)
	helpingErrorCallbackSending(api, upd, kb, out, ctx)
}

func HandleGroupSelection(api *maxbot.Api, upd *schemes.MessageCallbackUpdate, manager *postgres.Manager, ctx context.Context, group string) {
	intID, _ := strconv.Atoi(group)

	err := manager.SetUserGroup(upd.Callback.User.UserId, intID)
	if err != nil {
		logger.Logger.Printf("Error setting user group: %v", err)
		return
	}

	out := "Поздравляю вы выбрали группу блаблабла..."
	kb := keyboard.CreateKeyboardForStart(api)
	helpingErrorCallbackSending(api, upd, kb, out, ctx)
}

func helpingErrorCallbackSending(api *maxbot.Api, upd *schemes.MessageCallbackUpdate, kb *maxbot.Keyboard, out string, ctx context.Context) {
	res, err := api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetFormat("html").AddKeyboard(kb).SetText(out))
	if err != nil {
		logger.Logger.Printf("Error sending message: %v", err)
	} else {
		logger.Logger.Printf("Message sent: %v", res)
	}
}

func helpingErrorMessageSending(api *maxbot.Api, upd *schemes.MessageCallbackUpdate, out string, ctx context.Context) {
	res, err := api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetFormat("html").SetText(out))
	if err != nil {
		logger.Logger.Printf("Error sending message: %v", err)
	} else {
		logger.Logger.Printf("Message sent: %v", res)
	}
}
