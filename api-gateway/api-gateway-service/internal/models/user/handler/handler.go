package handler

import (
	"net/http"

	"strconv"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/user/grpc"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	log logger.Logger
	*grpc.UserGRPCClient
}

func NewUserHandler(log logger.Logger, cfg *config.Config) *UserHandler {
	grpcClient := grpc.NewUserGRPCClient(log, cfg) // Initialize the embedded type
	return &UserHandler{log: log, UserGRPCClient: grpcClient}
}

// @BasePath /api/v1/job
//
// @Summary Creates and schedule a job
// @Description Creates and schedule a job
// @Tags job
// @Accept json
// @Produce json
// @Success 200 {string} Job
// @Router /Create [get]
func (uh *UserHandler) CreateUser(c *gin.Context) {
	uh.log.Info("UserHandler.CreateUser: Entered")
	var createUserRequest models.CreateUserRequest
	if errValid := c.ShouldBindJSON(&createUserRequest); errValid != nil {
		uh.log.Errorf("UserHandler.CreateUser: error bind json with error: %V", errValid.Error())
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": errValid.Error()})
		return
	}

	res, err := uh.GRPC_CreateUser(c, createUserRequest.ToProto())

	if err != nil {
		uh.log.Errorf("UserHandler.CreateUser: User creation failed with error: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (uh *UserHandler) LoginUser(c *gin.Context) {
	uh.log.Info("UserHandler.LoginUser: Entered")

	var loginRequest models.LoginUserRequest
	if errValid := c.ShouldBindJSON(&loginRequest); errValid != nil {
		uh.log.Errorf("UserHandler.LoginUser: error bind json with error: %V", errValid.Error())
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": errValid})
		return
	}

	res, err := uh.GRPC_LoginUser(c, loginRequest.ToProto())

	if err != nil {
		uh.log.Errorf("UserHandler.LoginUser: User login failed with error: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  res.AccessToken,
		"refresh_token": res.RefreshToken,
	})
}

func (uh *UserHandler) RefreshAccessToken(c *gin.Context) {

	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": "Error getting refresh token."})
		return
	}

	res, err := uh.GRPC_RefreshAccessToken(c, &proto.RefreshTokenRequest{RefreshToken: refreshToken})

	if err != nil {
		uh.log.Errorf("UserHandler.RefreshAccessToken: refresh access token failed with error: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  res.AccessToken,
		"refresh_token": refreshToken,
	})
}

func (uh *UserHandler) LogOutUser(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	refreshToken := c.GetHeader("X-Refresh-Token")

	if refreshToken == "" || accessToken == "" {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": "Error getting tokens."})
		return
	}

	res, err := uh.GRPC_LogOutUser(c, &proto.LogOutRequest{AccessToken: accessToken, RefreshToken: refreshToken})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (uh *UserHandler) GetUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid UserId parameter"})
		return
	}

	res, err := uh.GRPC_GetUser(c, &proto.GetUserRequest{Id: int32(userId)})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err})
		return
	}

	c.JSON(http.StatusOK, res)
}
