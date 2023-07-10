package server

import (
	"fmt"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/docs"
	api "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/api"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	log logger.Logger
	cfg *config.Config
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {
	return &Server{log: log, cfg: cfg}
}

func (s *Server) Run() {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router := api.NewRouter(s.log, s.cfg)
	router.Setup(r)

	r.Run(fmt.Sprintf("0.0.0.0:%s", s.cfg.ServerPort))
}
