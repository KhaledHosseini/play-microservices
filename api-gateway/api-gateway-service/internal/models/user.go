package models

import (
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
)

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
	Role     int32
}

func (cur *CreateUserRequest) ToProto() *proto.CreateUserRequest {
	return &proto.CreateUserRequest{
		Name:     cur.Name,
		Email:    cur.Email,
		Password: cur.Password,
		Role:     proto.RoleType(cur.Role),
	}
}

type LoginUserRequest struct {
	Email    string
	Password string
}

func (cur *LoginUserRequest) ToProto() *proto.LoginUserRequest {
	return &proto.LoginUserRequest{
		Email:    cur.Email,
		Password: cur.Password,
	}
}

type RefreshTokenRequest struct {
	RefreshToken string
}

func (rtr *RefreshTokenRequest) ToProto() *proto.RefreshTokenRequest {
	return &proto.RefreshTokenRequest{
		RefreshToken: rtr.RefreshToken,
	}
}

type LogOutRequest struct {
	RefreshToken string
	AccessToken  string
}

func (lor *LogOutRequest) ToProto() *proto.LogOutRequest {
	return &proto.LogOutRequest{
		RefreshToken: lor.RefreshToken,
		AccessToken:  lor.RefreshToken,
	}
}
