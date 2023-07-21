package handler

import (
	"net/http"

	"strconv"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/user/grpc"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/cookie"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	log logger.Logger
	*grpc.UserGRPCClient
	cookieManager *cookie.CookieManager
}

func NewUserHandler(log logger.Logger, cfg *config.Config) *UserHandler {
	grpcClient := grpc.NewUserGRPCClient(log, cfg) // Initialize the embedded type
	cookieManager := cookie.NewCookie(cfg)
	return &UserHandler{log: log, UserGRPCClient: grpcClient, cookieManager: cookieManager}
}

// @Summary Register a user
// @Description Register a user
// @Tags user
// @Accept json
// @Produce json
// @Param   input      body    models.CreateUserRequest     true        "input"
// @Success 200 {object} models.CreateUserResponse
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
// @Success      200      {object}  models.LoginUserResponse
// @Header       200      {string}  Authorization     "Access-Token"
// @Header       200      {string}  X-Refresh-Token    "Refresh-Token"
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

	uh.cookieManager.SetAccessToken(c, res.AccessToken, res.AccessTokenAge)
	uh.cookieManager.SetRefreshToken(c, res.RefreshToken, res.RefreshTokenAge)

	c.JSON(http.StatusOK, &models.LoginUserResponse{Message: "Login success"})
}

// @Summary refresh access token
// @Description refresh access token
// @Tags user
// @Param   X-Refresh-Token      header    string     true        "X-Refresh-Token"
// @Success      200      {object}  models.RefreshTokenResponse
// @Header       200      {string}  Authorization     "Access-Token"
// @Header       200      {string}  X-Refresh-Token    "Refresh-Token"
// @Router /user/refresh_token [post]
func (uh *UserHandler) RefreshAccessToken(c *gin.Context) {

	uh.log.Info("UserHandler.RefreshAccessToken started")
	refreshToken, err := cookie.GetRefreshToken(c)
	if err != nil {
		uh.log.Info("UserHandler.RefreshAccessToken error getting refresh token")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "error": "Error getting refresh token."})
		return
	}

	res, err := uh.GRPC_RefreshAccessToken(c, &proto.RefreshTokenRequest{RefreshToken: refreshToken})

	if err != nil {
		uh.log.Errorf("UserHandler.RefreshAccessToken: refresh access token failed with error: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}
	uh.log.Info("UserHandler.RefreshAccessToken refresh token success.")

	uh.cookieManager.SetAccessToken(c, res.AccessToken, res.AccessTokenAge)

	c.JSON(http.StatusOK, &models.RefreshTokenResponse{Message: "access-token refresh success."})
}

// @Summary logout user
// @Description logout user
// @Tags user
// @Param   Authorization      header    string     true        "Authorization: Access token"
// @Param   X-Refresh-Token      header    string     true        "X-Refresh-Token"
// @Success      200      {object}  models.LogOutResponse
// @Router /user/logout [post]
func (uh *UserHandler) LogOutUser(c *gin.Context) {
	accessToken, err := cookie.GetAccessToken(c)
	if err != nil {
		uh.log.Info("UserHandler.LogOutUser error getting access token")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "error": "Error getting access token."})
		return
	}

	refreshToken, err2 := cookie.GetRefreshToken(c)
	if err2 != nil {
		uh.log.Info("UserHandler.LogOutUser error getting refresh token")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "error": "Error getting refresh token."})
		return
	}

	res, err := uh.GRPC_LogOutUser(c, &proto.LogOutRequest{AccessToken: accessToken, RefreshToken: refreshToken})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}
	uh.cookieManager.SetAccessToken(c, "", 0)
	uh.cookieManager.SetRefreshToken(c, "", 0)

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

	// get user id from cookie?

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

// @Summary Get the list of users
// @Description retrieve the list of users
// @Tags user
// @Produce json
// @Param   page      query    int     true        "Page"
// @Param   size      query    int     true        "Size"
// @Success 200 {array} models.ListUsersResponse
// @Router /user/list [get]
func (uh *UserHandler) ListUsers(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if err != nil {
		// Handle missing or invalid page parameter
		c.JSON(400, gin.H{"error": "Invalid page parameter"})
		return
	}

	size, err := strconv.ParseInt(c.Query("size"), 10, 32)
	if err != nil {
		// Handle missing or invalid size parameter
		c.JSON(400, gin.H{"error": "Invalid size parameter"})
		return
	}

	res, err := uh.GRPC_ListUsers(c, &proto.ListUsersRequest{Page: page, Size: size})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ListUsersResponseFromProto(res))
}
