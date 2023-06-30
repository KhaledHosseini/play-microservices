package server

import (
	"context"
	"log"
	"net"

	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/config"
	MyJobGRPCService "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models/job/grpc"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/logger"
	JobGRPCServiceProto "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/proto"

	JobDB "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models/job/database"
	kafka "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models/job/message_broker"

	"github.com/reugn/go-quartz/quartz"
	kg "github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	log       logger.Logger
	cfg       *config.Config
	mongoDB   *mongo.Client
	kafkaConn *kg.Conn
}

// NewServer constructor
func NewServer(log logger.Logger, cfg *config.Config, mongoDB *mongo.Client, kafkaConn *kg.Conn) *server {
	return &server{log: log, cfg: cfg, mongoDB: mongoDB, kafkaConn: kafkaConn}
}

func (s *server) Run() error {
	lis, err := net.Listen("tcp", ":"+s.cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sched := quartz.NewStdScheduler()
	sched.Start(ctx)
	defer func() {
		sched.Stop()
		sched.Wait(ctx)
	}()

	grpc_server := grpc.NewServer()
	job_db := JobDB.NewJobDBMongo(s.mongoDB)

	job_producer := kafka.NewJobsProducer()
	job_producer.Run(s.kafkaConn, s.cfg)
	defer job_producer.Close()

	job_service := MyJobGRPCService.NewJobService(s.log, job_db, job_producer, sched)
	job_service.LoadScheduledJobs()
	JobGRPCServiceProto.RegisterJobServiceServer(grpc_server, job_service)
	reflection.Register(grpc_server)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpc_server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}
