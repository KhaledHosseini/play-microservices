package models

import (
	"time"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/thoas/go-funk"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Job struct {
	Id           string
	Name         string
	Description  string
	ScheduleTime time.Time
	CreatedTime  time.Time
	UpdatedTime  time.Time
	JobStatus    int32
	JobType      int32
	JobData      string
}

func JobFromProto(p *proto.Job) *Job {
	return &Job{
		Id:           p.Id,
		Name:         p.Name,
		Description:  p.Description,
		ScheduleTime: p.ScheduleTime.AsTime(),
		CreatedTime:  p.CreatedTime.AsTime(),
		UpdatedTime:  p.UpdatedTime.AsTime(),
		JobStatus:    int32(p.Status),
		JobType:      int32(p.JobType),
		JobData:      p.JobData,
	}
}

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

type CreateJobResponse struct {
	Id string
}

func CreateJobResponseFromProto(p *proto.CreateJobResponse) *CreateJobResponse {
	return &CreateJobResponse{
		Id: p.Id,
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

type ListJobsResponse struct {
	TotalCount int64
	TotalPages int64
	Page       int64
	Size       int64
	HasMore    bool
	Jobs       []Job
}

func ListJobsResponseFromProto(p *proto.ListJobsResponse) *ListJobsResponse {

	jobs := funk.Map(p.Jobs, func(x *proto.Job) Job {
		return *JobFromProto(x)
	}).([]Job)

	return &ListJobsResponse{
		TotalCount: p.TotalCount,
		TotalPages: p.TotalPages,
		Page:       p.Page,
		Size:       p.Size,
		HasMore:    p.HasMore,
		Jobs:       jobs,
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

type UpdateJobResponse struct {
	Message string
}

func UpdateJobResponseFromProto(p *proto.UpdateJobResponse) *UpdateJobResponse {
	return &UpdateJobResponse{
		Message: p.Message,
	}
}

type DeleteJobResponse struct {
	Message string
}

func DeleteJobResponseFromProto(p *proto.DeleteJobResponse) DeleteJobResponse {
	return DeleteJobResponse{
		Message: p.Message,
	}
}
