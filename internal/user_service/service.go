package user_service

import (
	"context"
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
	inGroup, err := telegram.IsUserInGroup(req.ChannelId, req.UserId)
	if err != nil {
		return nil, err
	}

	return &proto.CheckUserInGroupResponse{
		IsUserInGroup: inGroup,
	}, nil
}
