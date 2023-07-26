package models

import (
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/thoas/go-funk"
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

type CreateUserResponse struct {
	Message string
}

func CreateUserResponseFromProto(p *proto.CreateUserResponse) *CreateUserResponse {
	return &CreateUserResponse{
		Message: p.Message,
	}
}

type LoginUserRequest struct {
	Email    string
	Password string
}

type LoginUserResponse struct {
	Message string
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

type RefreshTokenResponse struct {
	Message string
}

type LogOutRequest struct {
	RefreshToken string
	AccessToken  string
}

func (lor *LogOutRequest) ToProto() *proto.LogOutRequest {
	return &proto.LogOutRequest{
		RefreshToken: lor.RefreshToken,
	}
}

type LogOutResponse struct {
	Message string
}

type GetUserResponse struct {
	Id    int32
	Name  string
	Email string
}

func GetUserResponseFromProto(p *proto.GetUserResponse) *GetUserResponse {
	return &GetUserResponse{
		Id:    p.Id,
		Name:  p.Name,
		Email: p.Email,
	}
}

type ListUsersResponse struct {
	TotalCount int64
	TotalPages int64
	Page       int64
	Size       int64
	HasMore    bool
	Users      []GetUserResponse
}

func ListUsersResponseFromProto(p *proto.ListUsersResponse) *ListUsersResponse {
	users := funk.Map(p.Users, func(x *proto.GetUserResponse) GetUserResponse {
		return *GetUserResponseFromProto(x)
	}).([]GetUserResponse)

	return &ListUsersResponse{
		TotalCount: p.TotalCount,
		TotalPages: p.TotalPages,
		Page:       p.Page,
		Size:       p.Size,
		HasMore:    p.HasMore,
		Users:      users,
	}
}
