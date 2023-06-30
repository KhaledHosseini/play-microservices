package grpc

import (
	context "context"

	pb "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/proto"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type JobService struct {
	pb.UnimplementedJobServiceServer
}

func NewJobService() *JobService {
	return &JobService{}
}

func (j *JobService) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateJob not implemented")
}

func (JobService) GetJob(context.Context, *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJob not implemented")
}

func (JobService) ListJobs(context.Context, *pb.ListJobsRequest) (*pb.ListJobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListJobs not implemented")
}

func (JobService) UpdateJob(context.Context, *pb.UpdateJobRequest) (*pb.UpdateJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateJob not implemented")
}

func (JobService) DeleteJob(context.Context, *pb.DeleteJobRequest) (*pb.DeleteJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteJob not implemented")
}
