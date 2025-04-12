package user_service

import (
	"context"
	"fmt"
	"user_service/internal/telegram"
	"user_service/pkg/proto"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CheckUserInGroup(ctx context.Context, req *proto.CheckUserInGroupRequest) (*proto.CheckUserInGroupResponse, error) {
	if err := telegram.Init(req.Token); err != nil {
		return nil, fmt.Errorf("не удалось инициализировать Telegram: %w", err)
	}

	inGroup, err := telegram.IsUserInGroup(req.ChannelId, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке пользователя: %w", err)
	}

	return &proto.CheckUserInGroupResponse{
		IsUserInGroup: inGroup,
	}, nil
}
