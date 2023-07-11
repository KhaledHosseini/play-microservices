package api

import (
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/docs"
	interceptors "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/api/interceptors"
	jh "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/job/handler"
	rh "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/report/handler"
	uh "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/user/handler"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/gin-gonic/gin"

	"net/http"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	log logger.Logger
	cfg *config.Config
}

func NewRouter(log logger.Logger, cfg *config.Config) *Router {
	return &Router{log: log, cfg: cfg}
}
func (r *Router) Setup(router *gin.Engine) {

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("/api/v1/ping", r.Pong)

	userHandler := uh.NewUserHandler(r.log, r.cfg)
	jobHandler := jh.NewJobHandler(r.log, r.cfg)
	reportHandler := rh.NewReportHandler(r.log, r.cfg)

	router.POST("/api/v1/user/create", userHandler.CreateUser)
	router.POST("/api/v1/user/login", userHandler.LoginUser)

	router.Group("/api", interceptors.AuthenticateUser())
	{
		router.POST("/api/v1/user/refresh_token", userHandler.RefreshAccessToken)
		router.POST("/api/v1/user/logout", userHandler.LogOutUser)
		router.GET("/api/v1/user/:id", userHandler.GetUser)
		router.GET("/api/v1/user/list", userHandler.ListUsers)

		router.POST("/api/v1/job/create", jobHandler.CreateJob)
		router.GET("/api/v1/job/:id", jobHandler.GetJob)
		router.GET("/api/v1/job/list", jobHandler.ListJobs)
		router.POST("/api/v1/job/update", jobHandler.UpdateJob)
		router.DELETE("/api/v1/job/:id", jobHandler.DeleteJob)

		router.GET("/api/v1/report/list", reportHandler.ListReports)
	}
}

// @BasePath /api/v1
// PingExample godoc
// @Summary ping
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
