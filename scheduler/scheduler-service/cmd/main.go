package main

import (
    "log"

    "github.com/<your_username>/play-microservices/scheduler/scheduler-service/config"
"github.com/<your_username>/play-microservices/scheduler/scheduler-service/pkg/logger"
    "github.com/<your_username>/play-microservices/scheduler/scheduler-service/internal/server"
)

func main() {

    cfg, err := config.InitConfig()
    if err != nil {
        log.Fatal(err)
    }

    appLogger := logger.NewApiLogger(cfg)
    appLogger.InitLogger()
    appLogger.Info("Starting user server")
    appLogger.Infof(
        "AppVersion: %s, LogLevel: %s, Environment: %s",
        cfg.AppVersion,
        cfg.Logger_Level,
        cfg.Environment,
    )
    appLogger.Infof("Success parsed config: %#v", cfg.AppVersion)

    s := server.NewServer(appLogger, cfg)
    s.Run()
}