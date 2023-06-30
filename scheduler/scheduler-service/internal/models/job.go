package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/utils"
	jobsService "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Job struct {
	JobID        primitive.ObjectID `json:"jobId" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description  string             `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	ScheduleTime time.Time          `json:"scheduleTime" bson:"scheduleTime,omitempty"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	Status       int32              `json:"status,omitempty" bson:"status,omitempty"`
	JobType      int32              `json:"jobType,omitempty" bson:"jobType,omitempty"`
	JobData      string             `json:"jobData,omitempty" bson:"jobData,omitempty"`
	ScheduledKey int                `json:"scheduledKey,omitempty" bson:"scheduledKey,omitempty"`
}

func (j *Job) ToProto() *jobsService.Job {
	return &jobsService.Job{
		Id:           j.JobID.String(),
		Name:         j.Name,
		Description:  j.Description,
		ScheduleTime: timestamppb.New(j.ScheduleTime),
		CreatedTime:  timestamppb.New(j.CreatedAt),
		UpdatedTime:  timestamppb.New(j.UpdatedAt),
		Status:       jobsService.JobStatus(j.Status),
		JobType:      jobsService.JobType(j.JobType),
		JobData:      j.JobData,
	}
}

func JobFromProto_Job(jb *jobsService.Job) (*Job, error) {
	jobID, err := primitive.ObjectIDFromHex(jb.GetId())
	if err != nil {
		return nil, err
	}
	return &Job{
		JobID:        jobID,
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

func JobFromProto_CreateJobRequest(jbr *jobsService.CreateJobRequest) *Job {
	return &Job{
		Name:         jbr.GetName(),
		Description:  jbr.GetDescription(),
		ScheduleTime: jbr.GetScheduleTime().AsTime(),
		JobType:      int32(jbr.GetJobType()),
		JobData:      jbr.GetJobData(),
	}
}

type JobsList struct {
	TotalCount int64  `json:"totalCount"`
	TotalPages int64  `json:"totalPages"`
	Page       int64  `json:"page"`
	Size       int64  `json:"size"`
	HasMore    bool   `json:"hasMore"`
	Jobs       []*Job `json:"jobs"`
}

// databas interface for Job model
type JobDB interface {
	Create(ctx context.Context, job *Job) (*Job, error)
	Update(ctx context.Context, job *Job) (*Job, error)
	GetByID(ctx context.Context, jobID primitive.ObjectID) (*Job, error)
	DeleteByID(ctx context.Context, jobID primitive.ObjectID) error
	GetByScheduledKey(ctx context.Context, jobScheduledKey int) (*Job, error)
	DeleteByScheduledKey(ctx context.Context, jobScheduledKey int) error
	ListALL(ctx context.Context, pagination *utils.Pagination) (*JobsList, error)
}

// Message broker interface for Job model
type JobsProducer interface {
	PublishCreate(ctx context.Context, job *Job) error
	PublishUpdate(ctx context.Context, job *Job) error
	PublishRun(ctx context.Context, job *Job) error
}
