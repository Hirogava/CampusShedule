package maxbot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hirogava/CampusShedule/internal/config/logger"
	botHanlders "github.com/Hirogava/CampusShedule/internal/maxbot"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/configservice"
)

func SetMaxConf() {
	var configPath string
	var maxenv = os.Getenv("MAXBOT_ENV")
	// Customize ConsoleWriter
	
	if maxenv != "" {
		configPath = "/go/bin/config/app-" + maxenv + ".yaml"
	} else if 2 <= len(os.Args) {
		configPath = os.Args[1]
	} else {
		logger.Logger.Error("maxenv environment variable not found. Stop.")
		return
	}

	configService := configservice.NewConfigInterface(configPath)
	if configService == nil {
		logger.Logger.Fatal("configPath", configPath, "NewConfigInterface failed. Stop.")
	}

	api, err := maxbot.NewWithConfig(configService)
	if err != nil {
		logger.Logger.Fatal("NewWithConfig failed. Stop.", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGTERM, os.Interrupt)
		<-exit
		cancel()
	}()

	info, err := api.Bots.GetBot(ctx) // Простой метод
	logger.Logger.Printf("Get me: %#v %#v", info, err)

	chatList, err := api.Chats.GetChats(ctx, 0, 0)
	if err != nil {
		fmt.Printf("Unknown type: %#v", err)
	}
	for _, chat := range chatList.Chats {
		fmt.Printf("Bot is members at the chat: %#v", chat.Title)
		fmt.Printf("	: %#v", chat.ChatId)
	}

	botHanlders.StartListening(api, ctx)
}
