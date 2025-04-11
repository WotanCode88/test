package user_service

import (
	"fmt"
	"user_service/internal/db"
	"user_service/internal/telegram"
)

func CheckAndAddUserToDB(token string, channelID int64, userID int64) error {
	isMember, err := telegram.IsUserInGroup(channelID, userID)
	if err != nil {
		return fmt.Errorf("ошибка: %v", err)
	}

	if !isMember {
		return fmt.Errorf("юзер не в группе")
	}

	query := `INSERT INTO users (telegram_id, channel_id) VALUES ($1, $2) RETURNING id`
	_, err = db.DB.Exec(query, userID, channelID)
	if err != nil {
		return fmt.Errorf("ошибка бд: %v", err)
	}

	return nil
}
