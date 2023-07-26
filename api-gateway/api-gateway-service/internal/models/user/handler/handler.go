package handler

import (
	"net/http"

	"strconv"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/user/grpc"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/cookie"
	grpcutils "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/grpc"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
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
		uh.log.Errorf("UserHandler.CreateUser: error get input with error: %v", errValid.Error())
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": errValid.Error()})
		return
	}

	res, err := uh.GRPC_CreateUser(c, createUserRequest.ToProto())
	if err != nil {
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		uh.log.Errorf("UserHandler.CreateUser: User creation failed with error: %v", err.Error())
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
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
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"status": false, "error": errValid})
		return
	}

	res, err := uh.GRPC_LoginUser(c, loginRequest.ToProto())

	if err != nil {
		uh.log.Errorf("UserHandler.LoginUser: User login failed with error: %v", err.Error())
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
		return
	}

	uh.cookieManager.SetAccessToken(c, res.AccessToken, res.AccessTokenAge)
	uh.cookieManager.SetRefreshToken(c, res.RefreshToken, res.RefreshTokenAge)

	c.JSON(http.StatusOK, &models.LoginUserResponse{Message: "Login success"})
}

// @Summary refresh access token
// @Description refresh access token
// @Tags user
// @Success      200      {object}  models.RefreshTokenResponse
// @Header       200      {string}  Authorization     "Access-Token"
// @Router /user/refresh_token [post]
func (uh *UserHandler) RefreshAccessToken(c *gin.Context) {

	uh.log.Info("UserHandler.RefreshAccessToken started")
	refreshToken, err := cookie.GetRefreshToken(c)
	if err != nil {
		uh.log.Errorf("UserHandler.RefreshAccessToken: error getting refresh cookie %v", err.Error())
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "error": "Error getting refresh token."})
		return
	}

	res, err := uh.GRPC_RefreshAccessToken(c, &proto.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		uh.log.Errorf("UserHandler.RefreshAccessToken: error refreshing cookie: %v", err.Error())
		//If the error is unauthenticated, it means our refresh token has expired. so we clear the auth cookies.
		grpcCode := grpcutils.GetGRPCStatus(err).Code()
		switch grpcCode {
		case codes.Unauthenticated:
			uh.cookieManager.SetAccessToken(c, "", 0)
			uh.cookieManager.SetRefreshToken(c, "", 0)
			c.JSON(http.StatusNetworkAuthenticationRequired, &models.RefreshTokenResponse{Message: "please login again."})
			return
		}
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
		return
	}
	uh.log.Info("UserHandler.RefreshAccessToken refresh token success.")
	uh.cookieManager.SetAccessToken(c, res.AccessToken, res.AccessTokenAge)

	c.JSON(http.StatusOK, &models.RefreshTokenResponse{Message: "access-token refresh success."})
}

// @Summary logout user
// @Description logout user
// @Tags user
// @Success      200      {object}  models.LogOutResponse
// @Router /user/logout [post]
func (uh *UserHandler) LogOutUser(c *gin.Context) {

	refreshToken, err := cookie.GetRefreshToken(c)
	if err != nil {
		uh.log.Errorf("UserHandler.LogOutUser: error getting refresh cookie %v", err.Error())
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "error": "Error getting refresh token."})
		return
	}

	res, err := uh.GRPC_LogOutUser(c, &proto.LogOutRequest{RefreshToken: refreshToken})
	if err != nil {
		uh.log.Errorf("UserHandler.LogOutUser error logging out: %v", err.Error())
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
		return
	}
	uh.cookieManager.SetAccessToken(c, "", 0)
	uh.cookieManager.SetRefreshToken(c, "", 0)

	c.JSON(http.StatusOK, &models.LogOutResponse{Message: res.Message})
}

// @Summary Get user
// @Description Get user
// @Tags user
// @Produce json
// @Success 200 {object} models.GetUserResponse
// @Router /user/get [get]
func (uh *UserHandler) GetUser(c *gin.Context) {

	uh.log.Infof("Enter GetUser: %v", c.Request.Context())

	res, err := uh.GRPC_GetUser(c.Request.Context(), &proto.GetUserRequest{})
	if err != nil {
		uh.log.Errorf("UserHandler.GetUser error getting user: %v", err.Error())
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
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
		uh.log.Errorf("UserHandler.ListUsers no page parameters is provided: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	size, err := strconv.ParseInt(c.Query("size"), 10, 32)
	if err != nil {
		uh.log.Errorf("UserHandler.ListUsers no size parameters is provided: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
		return
	}

	res, err := uh.GRPC_ListUsers(c.Request.Context(), &proto.ListUsersRequest{Page: page, Size: size})
	if err != nil {
		uh.log.Errorf("UserHandler.ListUsers error getting the user: %v", err.Error())
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ListUsersResponseFromProto(res))
}
