package server

import (
	"fmt"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	api "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/api"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"

	"github.com/gin-gonic/gin"
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

	router := api.NewRouter(s.log, s.cfg)
	router.Setup(r)

	r.Run(fmt.Sprintf("0.0.0.0:%s", s.cfg.ServerPort))
}
