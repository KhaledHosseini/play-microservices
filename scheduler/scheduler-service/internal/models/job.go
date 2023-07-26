package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	proto "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/proto"
	"github.com/thoas/go-funk"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Job struct {
	Id        primitive.ObjectID 	`json:"Id" bson:"_id,omitempty"`
	Name         string             `json:"Name" bson:"name" validate:"required,min=3,max=250"`
	Description  string             `json:"Description" bson:"description" validate:"required,min=3,max=500"`
	ScheduleTime time.Time          `json:"ScheduleTime" bson:"scheduleTime"`
	CreatedAt    time.Time          `json:"CreatedAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"UpdatedAt" bson:"updatedAt"`
	Status       int32              `json:"Status" bson:"status"`
	JobType      int32              `json:"JobType" bson:"jobType"`
	JobData      string             `json:"JobData" bson:"jobData"`
	ScheduledKey int                `json:"ScheduledKey" bson:"scheduledKey"`
}

func (j *Job) ToProto() *proto.Job {
	return &proto.Job{
		Id:           j.Id.Hex(),
		Name:         j.Name,
		Description:  j.Description,
		ScheduleTime: timestamppb.New(j.ScheduleTime),
		CreatedTime:  timestamppb.New(j.CreatedAt),
		UpdatedTime:  timestamppb.New(j.UpdatedAt),
		Status:       proto.JobStatus(j.Status),
		JobType:      proto.JobType(j.JobType),
		JobData:      j.JobData,
	}
}

// Note: 0 values for int and empty strings will be ignored proto.
func JobFromProto_Job(jb *proto.Job) (*Job, error) {
	jobId, err := primitive.ObjectIDFromHex(jb.GetId())
	if err != nil {
		return nil, err
	}
	return &Job{
		Id:        	  jobId,
		Name:         jb.GetName(),
		Description:  jb.GetDescription(),
		ScheduleTime: jb.GetScheduleTime().AsTime(),
		CreatedAt:    jb.GetCreatedTime().AsTime(),
		UpdatedAt:    jb.GetUpdatedTime().AsTime(),
		Status:       int32(jb.GetStatus()),
		JobType:      int32(jb.GetJobType()),
		JobData:      jb.GetJobData(),
		ScheduledKey: 0,
	}, nil
}

func JobFromProto_CreateJobRequest(jbr *proto.CreateJobRequest) (*Job, error) {
	return &Job{
		Name:         jbr.GetName(),
		Description:  jbr.GetDescription(),
		ScheduleTime: jbr.GetScheduleTime().AsTime(),
		JobType:      int32(jbr.GetJobType()),
		JobData:      jbr.GetJobData(),
	}, nil
}

func JobFromProto_UpdateJobRequest(jbr *proto.UpdateJobRequest) (*Job, error) {
	jobId, err := primitive.ObjectIDFromHex(jbr.GetId())
	if err != nil {
		return nil, err
	}
	return &Job{
		Id:           jobId,
		Name:         jbr.GetName(),
		Description:  jbr.GetDescription(),
		ScheduleTime: jbr.GetScheduleTime().AsTime(),
		JobType:      int32(jbr.GetJobType()),
		JobData:      jbr.GetJobData(),
	}, nil
}

type JobsList struct {
	TotalCount int64  `json:"totalCount"`
	TotalPages int64  `json:"totalPages"`
	Page       int64  `json:"page"`
	Size       int64  `json:"size"`
	HasMore    bool   `json:"hasMore"`
	Jobs       []*Job `json:"jobs"`
}

func (j *JobsList) ToProto() *proto.ListJobsResponse {
	jobs := funk.Map(j.Jobs, func(x *Job) *proto.Job {
		return x.ToProto()
	}).([]*proto.Job)
	return &proto.ListJobsResponse{
		TotalCount: j.TotalCount,
		TotalPages: j.TotalPages,
		Page:       j.Page,
		Size:       j.Size,
		HasMore:    j.HasMore,
		Jobs:       jobs,
	}
}

// databas interface for Job model
type JobDB interface {
	Create(ctx context.Context, job *Job) (*Job, error)
	Update(ctx context.Context, job *Job) (*Job, error)
	GetByID(ctx context.Context, jobID string) (*Job, error)
	DeleteByID(ctx context.Context, jobID string) error
	GetByScheduledKey(ctx context.Context, jobScheduledKey int) (*Job, error)
	DeleteByScheduledKey(ctx context.Context, jobScheduledKey int) error
	ListALL(ctx context.Context, page int64, size int64) (*JobsList, error)
}

// Message broker interface for Job model
type JobsProducer interface {
	PublishCreate(ctx context.Context, job *Job) error
	PublishUpdate(ctx context.Context, job *Job) error
	PublishRun(ctx context.Context, job *Job) error
}
