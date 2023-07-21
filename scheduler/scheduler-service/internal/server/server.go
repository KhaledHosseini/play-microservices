package server

import (
	"context"
	"net"

	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/config"
	Interceptors "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/interceptors"
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
		s.log.Fatalf("failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.log.Info("Starting quartz scheduler...")
	sched := quartz.NewStdScheduler()
	sched.Start(ctx)
	defer func() {
		s.log.Info("Stoping quartz scheduler...")
		sched.Stop()
		sched.Wait(ctx)
	}()

	job_producer := kafka.NewJobsProducer(s.log)
	job_producer.Run(s.kafkaConn, s.cfg)
	defer job_producer.Close()
	
	auth_interceptor := Interceptors.NewAuthInterceptor(s.log,s.cfg)
	grpc_server := grpc.NewServer(
		grpc.UnaryInterceptor(auth_interceptor.AuthInterceptor),
	)
	job_db := JobDB.NewJobDBMongo(s.mongoDB)

	JobService := MyJobGRPCService.NewJobService(s.log, job_db, job_producer, sched)
	JobService.LoadScheduledJobs()
	JobGRPCServiceProto.RegisterJobsServiceServer(grpc_server, JobService)
	reflection.Register(grpc_server)

	s.log.Info("Starting kafka consumer...")
	jobsConsumer := kafka.NewJobsConsumerGroup(s.log, job_db)
	jobsConsumer.Run(ctx, cancel, s.kafkaConn, s.cfg)

	s.log.Infof("gRPC server listening at %v", lis.Addr())
	if err := grpc_server.Serve(lis); err != nil {
		s.log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}
