package telegram

import (
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	bot      *tgbotapi.BotAPI
	initMux  sync.Mutex
	botToken string
)

func Init(token string) error {
	initMux.Lock()
	defer initMux.Unlock()

	if bot != nil && botToken == token {
		return nil
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("ошибка авторизации бота: %v", err)
	}

	botToken = token
	log.Printf("бот авторизован как: %s", bot.Self.UserName)
	return nil
}

func IsUserInGroup(chatID int64, userID int64) (bool, error) {
	if bot == nil {
		return false, fmt.Errorf("бот не инициализирован — вызови Init(token)")
	}

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
