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

// @Summary Register a user
// @Description Register a user
// @Tags user
// @Accept json
// @Produce json
// @Param   input      body    models.CreateUserRequest     true        "input"
// @Success 200 {object} models.CreateJobResponse
// @Router /user/create [post]
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

	c.JSON(http.StatusOK, models.CreateUserResponseFromProto(res))
}

// @Summary login a user
// @Description login a user
// @Tags user
// @Accept json
// @Param   input      body    models.LoginUserRequest     true        "input"
// @Success      200      {string}  string    "success"
// @Header       200      {string}  Access-Token     "Access-Token"
// @Header       200      {string}  Refresh-Token    "Refresh-Token"
// @Router /user/login [post]
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.Header("Authorization", res.AccessToken)
	c.Header("X-Refresh-Token", res.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "login Success",
	})
}

// @Summary refresh access token
// @Description refresh access token
// @Tags user
// @Param   X-Refresh-Token      header    string     true        "X-Refresh-Token"
// @Success      200      {string}  string    "success"
// @Header       200      {string}  Access-Token     "Access-Token"
// @Header       200      {string}  Refresh-Token    "Refresh-Token"
// @Router /user/refresh_token [post]
func (uh *UserHandler) RefreshAccessToken(c *gin.Context) {

	uh.log.Info("UserHandler.RefreshAccessToken started")
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		uh.log.Info("UserHandler.RefreshAccessToken error getting refresh token")
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": "Error getting refresh token."})
		return
	}

	res, err := uh.GRPC_RefreshAccessToken(c, &proto.RefreshTokenRequest{RefreshToken: refreshToken})

	if err != nil {
		uh.log.Errorf("UserHandler.RefreshAccessToken: refresh access token failed with error: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}
	uh.log.Info("UserHandler.RefreshAccessToken refresh token success.")
	c.Header("Authorization", res.AccessToken)
	c.Header("X-Refresh-Token", refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

// @Summary logout user
// @Description logout user
// @Tags user
// @Param   Authorization      header    string     true        "Authorization: Access token"
// @Param   X-Refresh-Token      header    string     true        "X-Refresh-Token"
// @Success      200      {object}  models.LogOutResponse
// @Router /user/logout [post]
func (uh *UserHandler) LogOutUser(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	refreshToken := c.GetHeader("X-Refresh-Token")

	if refreshToken == "" || accessToken == "" {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": "Error getting tokens."})
		return
	}

	res, err := uh.GRPC_LogOutUser(c, &proto.LogOutRequest{AccessToken: accessToken, RefreshToken: refreshToken})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &models.LogOutResponse{Message: res.Message})
}

// @Summary Get user by id
// @Description Get user by id
// @Tags user
// @Produce json
// @Param   id      path    string     true        "some id"
// @Success 200 {object} models.GetUserResponse
// @Router /user/{id} [get]
func (uh *UserHandler) GetUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid or no id parameter"})
		return
	}

	res, err := uh.GRPC_GetUser(c, &proto.GetUserRequest{Id: int32(userId)})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.GetUserResponseFromProto(res))
}
