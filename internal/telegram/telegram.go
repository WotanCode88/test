package telegram

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI

func InitTelegram() error {
	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("TG_BOT_TOKEN не задан")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("ошибка авторизации бота: %v", err)
	}

	log.Printf("бот авторизован как: %s", bot.Self.UserName)
	return nil
}

func IsUserInGroup(chatID int64, userID int64) (bool, error) {
	config := tgbotapi.ChatConfigWithUser{
		ChatID: chatID,
		UserID: int(userID),
	}

	member, err := bot.GetChatMember(config)
	if err != nil {
		return false, fmt.Errorf("ошибка получения чата: %v", err)
	}

	switch member.Status {
	case "creator", "administrator", "member":
		return true, nil
	default:
		return false, nil
	}
}
