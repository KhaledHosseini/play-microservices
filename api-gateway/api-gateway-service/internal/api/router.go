package api

import (
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	interceptors "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/api/interceptors"
	jh "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/job/handler"
	rh "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/report/handler"
	uh "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/user/handler"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/gin-gonic/gin"

	"net/http"
)

type Router struct {
	log logger.Logger
	cfg *config.Config
}

func NewRouter(log logger.Logger, cfg *config.Config) *Router {
	return &Router{log: log, cfg: cfg}
}
func (r *Router) Setup(router *gin.Engine) {

	router.GET("/api/v1/ping", r.Pong)

	jobHandler := jh.NewJobHandler(r.log, r.cfg)
	reportHandler := rh.NewReportHandler(r.log, r.cfg)
	userHandler := uh.NewUserHandler(r.log, r.cfg)

	router.Group("/api", interceptors.AttachAccessTokenToGRPC())
	{
		router.POST("/api/v1/job/create", jobHandler.CreateJob)
		router.GET("/api/v1/job/get", jobHandler.GetJob)
		router.GET("/api/v1/job/list", jobHandler.ListJobs)
		router.POST("/api/v1/job/update", jobHandler.UpdateJob)
		router.POST("/api/v1/job/delete", jobHandler.DeleteJob)

		router.GET("/api/v1/report/list", reportHandler.ListReports)

		router.POST("/api/v1/user/create", userHandler.CreateUser)
		router.POST("/api/v1/user/login", userHandler.LoginUser)
		router.POST("/api/v1/user/refresh_token", userHandler.RefreshAccessToken)
		router.POST("/api/v1/user/logout", userHandler.LogOutUser)
		router.GET("/api/v1/user/get", userHandler.GetUser)
	}
}

// @BasePath /api/v1
// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {string} Pong
// @Router /ping [get]
func (s *Router) Pong(c *gin.Context) {
	c.JSON(http.StatusOK, "Pong")
}
