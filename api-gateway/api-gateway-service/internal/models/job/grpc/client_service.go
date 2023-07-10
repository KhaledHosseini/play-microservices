package grpc

import (
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	gr "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/grpc"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type JobGRPCClient struct {
	log logger.Logger
	cfg *config.Config
	gr.GRPC_Client
}

func NewJobGRPCCLient(log logger.Logger, cfg *config.Config) *JobGRPCClient {
	return &JobGRPCClient{log: log, cfg: cfg}
}

func (jc *JobGRPCClient) getClient() (proto.JobsServiceClient, *grpc.ClientConn, error) {
	conn, err := jc.Connect(jc.cfg.SchedulerServiceURL)
	if err != nil {
		return nil, nil, err
	}

	return proto.NewJobsServiceClient(conn), conn, nil
}

func (jc *JobGRPCClient) GRPC_CreateJob(c *gin.Context, createJobRequest *proto.CreateJobRequest) (*proto.CreateJobResponse, error) {
	client, conn, err := jc.getClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return client.CreateJob(c, createJobRequest)
}

func (jc *JobGRPCClient) GRPC_GetJob(c *gin.Context, getJobRequest *proto.GetJobRequest) (*proto.GetJobResponse, error) {

	client, conn, err := jc.getClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.GetJob(c, getJobRequest)
}

func (jc *JobGRPCClient) GRPC_ListJobs(c *gin.Context, listJobsRequest *proto.ListJobsRequest) (*proto.ListJobsResponse, error) {
	jc.log.Info("Connecting to grpc server.")
	client, conn, err := jc.getClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.ListJobs(c, listJobsRequest)
}

func (jc *JobGRPCClient) GRPC_UpdateJob(c *gin.Context, updateJobRequest *proto.UpdateJobRequest) (*proto.UpdateJobResponse, error) {

	client, conn, err := jc.getClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.UpdateJob(c, updateJobRequest)
}

func (jc *JobGRPCClient) GRPC_DeleteJob(c *gin.Context, deleteJobRequest *proto.DeleteJobRequest) (*proto.DeleteJobResponse, error) {

	client, conn, err := jc.getClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return client.DeleteJob(c, deleteJobRequest)
}
