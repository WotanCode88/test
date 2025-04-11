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
		return fmt.Errorf("такого бота не существует")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("не удалось создать: %v", err)
	}

	log.Printf("авторизация в %s", bot.Self.UserName)
	return nil
}

func IsUserInGroup(chatID int64, userID int64) (bool, error) {
	config := tgbotapi.ChatConfigWithUser{
		ChatID: chatID,
		UserID: int(userID),
	}

	member, err := bot.GetChatMember(config)
	if err != nil {
		return false, fmt.Errorf("ошибка с chatID: %v", err)
	}

	if member.Status == "member" {
		return true, nil // Пользователь в группе
	}

	return false, nil // Пользователь не в группе
}
