package grpc

import (
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	gr "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/grpc"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type UserGRPCClient struct {
	log logger.Logger
	cfg *config.Config
	gr.GRPC_Client
}

func NewUserGRPCClient(log logger.Logger, cfg *config.Config) *UserGRPCClient {
	return &UserGRPCClient{log: log, cfg: cfg}
}

func (uc *UserGRPCClient) getClient() (proto.UserServiceClient, *grpc.ClientConn, error) {
	uc.log.Info("Connecting to grpc server.")
	conn, err := uc.Connect(uc.cfg.AuthServiceURL)
	if err != nil {
		return nil, nil, err
	}
	uc.log.Info("Connected to grpc server.")
	return proto.NewUserServiceClient(conn), conn, nil
}

func (jc *UserGRPCClient) getClient2() (proto.UserProfileServiceClient, *grpc.ClientConn, error) {
	conn, err := jc.Connect(jc.cfg.AuthServiceURL)
	if err != nil {
		return nil, nil, err
	}
	return proto.NewUserProfileServiceClient(conn), conn, nil
}

func (uc *UserGRPCClient) GRPC_CreateUser(c *gin.Context, createUserRequest *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {

	client, conn, err := uc.getClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	uc.log.Info("Creating user.")
	return client.CreateUser(c, createUserRequest)
}

func (uc *UserGRPCClient) GRPC_LoginUser(c *gin.Context, loginUserRequest *proto.LoginUserRequest) (*proto.LoginUserResponse, error) {

	client, conn, err := uc.getClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.LoginUser(c, loginUserRequest)
}

func (uc *UserGRPCClient) GRPC_RefreshAccessToken(c *gin.Context, refreshTokenRequest *proto.RefreshTokenRequest) (*proto.RefreshTokenResponse, error) {

	client, conn, err := uc.getClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.RefreshAccessToken(c, refreshTokenRequest)
}

func (uc *UserGRPCClient) GRPC_LogOutUser(c *gin.Context, logOutRequest *proto.LogOutRequest) (*proto.LogOutResponse, error) {

	client, conn, err := uc.getClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.LogOutUser(c, logOutRequest)
}

func (uc *UserGRPCClient) GRPC_GetUser(c *gin.Context, getUserRequest *proto.GetUserRequest) (*proto.GetUserResponse, error) {

	client, conn, err := uc.getClient2()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.GetUser(c, getUserRequest)
}

func (uc *UserGRPCClient) GRPC_ListUsers(c *gin.Context, listUsersRequest *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {

	client, conn, err := uc.getClient2()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.ListUsers(c, listUsersRequest)
}
