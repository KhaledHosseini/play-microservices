package main

import (
	"log"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/server"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Environment: %s",
		cfg.AppVersion,
		cfg.Logger_Level,
		cfg.Environment,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.AppVersion)

	appLogger.Info("Starting the server")
	s := server.NewServer(appLogger, cfg)
	s.Run()
}
