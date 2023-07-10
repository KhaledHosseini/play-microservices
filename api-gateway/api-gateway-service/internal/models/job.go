package models

import (
	"time"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type CreateJobRequest struct {
	Name         string
	Description  string
	ScheduleTime time.Time
	JobType      int32
	JobData      string
}

func (cjr *CreateJobRequest) ToProto() *proto.CreateJobRequest {
	return &proto.CreateJobRequest{
		Name:         cjr.Name,
		Description:  cjr.Description,
		ScheduleTime: timestamppb.New(cjr.ScheduleTime),
		JobType:      proto.JobType(cjr.JobType),
		JobData:      cjr.JobData,
	}
}

type GetJobRequest struct {
	Id string
}

func (gjr *GetJobRequest) ToProto() *proto.GetJobRequest {
	return &proto.GetJobRequest{
		Id: gjr.Id,
	}
}

type ListJobsRequest struct {
	Page int64
	Size int64
}

func (ljr *ListJobsRequest) ToProto() *proto.ListJobsRequest {
	return &proto.ListJobsRequest{
		Page: ljr.Page,
		Size: ljr.Size,
	}
}

type UpdateJobRequest struct {
	Id           string
	Name         string
	Description  string
	ScheduleTime time.Time
	JobType      int32
	JobData      string
}

func (ujr *UpdateJobRequest) ToProto() *proto.UpdateJobRequest {
	return &proto.UpdateJobRequest{
		Id:           ujr.Id,
		Name:         ujr.Name,
		Description:  ujr.Description,
		ScheduleTime: timestamppb.New(ujr.ScheduleTime),
		JobType:      proto.JobType(ujr.JobType),
		JobData:      ujr.JobData,
	}
}

type DeleteJobRequest struct {
	Id string
}

func (gjr *DeleteJobRequest) ToProto() *proto.DeleteJobRequest {
	return &proto.DeleteJobRequest{
		Id: gjr.Id,
	}
}
