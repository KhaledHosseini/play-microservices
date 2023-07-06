package main

import (
	"context"
	"log"

	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/config"

	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/server"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/kafka"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/logger"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/mongodb"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	appLogger.Info("Connecting to mongodb server")
	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, cfg.DatabaseURI, cfg.DatabaseUser, cfg.DatabasePass)
	if err != nil {
		appLogger.Fatalf("cannot connect mongodb. uri is %s and error is %s", cfg.DatabaseURI, err)
	}
	defer func() {
		if err := mongoDBConn.Disconnect(ctx); err != nil {
			appLogger.Fatalf("mongoDBConn.Disconnect")
		}
	}()
	appLogger.Info("Connected to mongodb server")

	appLogger.Info("Connecting to kafka server")
	conn, err := kafka.NewKafkaConn(cfg.KafkaBrokers[0])
	if err != nil {
		appLogger.Fatalf("NewKafkaConn with error %s:", err)
	}
	defer conn.Close()

	brokers, err := conn.Brokers()
	if err != nil {
		appLogger.Fatalf("conn.Brokers with error %s", err)
	}
	appLogger.Infof("Kafka connected: %v", brokers)

	appLogger.Info("Starting the server")
	s := server.NewServer(appLogger, cfg, mongoDBConn, conn)
	s.Run()
}
