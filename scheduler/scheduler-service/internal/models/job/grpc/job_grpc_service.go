package grpc

import (
	context "context"
	"fmt"
	"math"
	"time"

	models "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models"
	validator "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models/job/validation"
	grpcErrors "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/grpc_errors"
	jobFucntion "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/job_function"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/logger"
	proto "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/proto"

	scheduler "github.com/reugn/go-quartz/quartz"
	"github.com/segmentio/fasthash/fnv1a"
)

type JobService struct {
	log           logger.Logger
	jobDB         models.JobDB
	jobsProducer  models.JobsProducer
	jobsScheduler scheduler.Scheduler
	proto.JobsServiceServer
}

func NewJobService(log logger.Logger, jobDB models.JobDB, jobsProducer models.JobsProducer, jobsScheduler scheduler.Scheduler) *JobService {
	return &JobService{log: log, jobDB: jobDB, jobsProducer: jobsProducer, jobsScheduler: jobsScheduler}
}

// exported method (Started with uppercase)
func (js *JobService) LoadScheduledJobs() {
	// load all pending jobs to the scheduler again
}
func (js *JobService) CreateJob(ctx context.Context, req *proto.CreateJobRequest) (*proto.CreateJobResponse, error) {

	v := validator.CreateJobRequestValidator{CreateJobRequest: req}
	v_err := v.Validate()
	if v_err != nil {
		return nil, grpcErrors.ErrInvalidInput
	}

	js.log.Infof("JobService.CreateJob: grpc message arrived : %v", req)
	job, err0 := models.JobFromProto_CreateJobRequest(req)
	if err0 != nil {
		return nil, grpcErrors.ErrorResponse(err0, err0.Error())
	}

	job.Status = int32(proto.JobStatus_JOB_STATUS_SCHEDULED)
	jobFingerPrint := fmt.Sprintf("%s:%s:%s:%p", job.Name, job.Description, job.JobData, &job.ScheduleTime)
	job.ScheduledKey = int(fnv1a.HashString64(jobFingerPrint)) //We assume 64 bit systems!
	//store the job in the database
	js.log.Infof("JobService.CreateJob: Creating job in the database")
	created, err := js.jobDB.Create(ctx, job)
	if err != nil {
		js.log.Errorf("JobService.CreateJob: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	jobID := created.JobID.Hex()
	functionJob := jobFucntion.NewFunctionJobWithKey(job.ScheduledKey, func(_ context.Context) (int, error) {
		js.log.Info("Job is triggered")
		jb, err := js.jobDB.GetByID(ctx, jobID)
		if err != nil {
			js.log.Errorf("JobService.JobFunction.GetByID: %v", err)
			return 0, grpcErrors.ErrorResponse(err, err.Error())
		}
		//set the state of the job to running
		jb.Status = int32(proto.JobStatus_JOB_STATUS_RUNNING)
		_, err2 := js.jobDB.Update(ctx, jb)
		if err2 != nil {
			js.log.Errorf("JobService.JobFunction.Update: %v", err)
			return 0, grpcErrors.ErrorResponse(err2, err2.Error())
		}

		js.jobsProducer.PublishRun(ctx, jb)

		return 0, nil
	})

	triggerTime := time.Duration(math.Max(0, time.Until(job.ScheduleTime).Seconds()))
	js.log.Infof("triggers in %s seconds", triggerTime)

	js.jobsScheduler.ScheduleJob(ctx, functionJob, scheduler.NewRunOnceTrigger(triggerTime))

	js.jobsProducer.PublishCreate(ctx, created)

	return &proto.CreateJobResponse{Id: created.JobID.String()}, nil
}

func (js *JobService) GetJob(ctx context.Context, req *proto.GetJobRequest) (*proto.Job, error) {
	js.log.Infof("JobService.GetJob: grpc message arrived : %v", req)

	v := validator.GetJobRequestValidator{GetJobRequest: req}
	v_err := v.Validate()
	if v_err != nil {
		return nil, grpcErrors.ErrInvalidInput
	}

	job, err := js.jobDB.GetByID(ctx, req.GetId())
	if err != nil {
		js.log.Errorf("JobService.GetJob: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	return job.ToProto(), nil
}

func (js *JobService) ListJobs(ctx context.Context, req *proto.ListJobsRequest) (*proto.ListJobsResponse, error) {
	js.log.Infof("JobService.ListJobs: grpc message arrived : %v", req)

	v := validator.ListJobsRequestValidator{ListJobsRequest: req}
	v_err := v.Validate()
	if v_err != nil {
		return nil, grpcErrors.ErrInvalidInput
	}

	jobs, err := js.jobDB.ListALL(ctx, req.Size, req.Page)
	if err != nil {
		js.log.Errorf("JobService.ListJobs: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	js.log.Infof("jobs returned: %v", jobs)

	return jobs.ToProto(), nil
}

func (js *JobService) UpdateJob(ctx context.Context, req *proto.UpdateJobRequest) (*proto.UpdateJobResponse, error) {
	js.log.Infof("JobService.UpdateJob: grpc message arrived : %v", req)

	v := validator.UpdateJobRequestValidator{UpdateJobRequest: req}
	v_err := v.Validate()
	if v_err != nil {
		return nil, grpcErrors.ErrInvalidInput
	}

	job, err := models.JobFromProto_UpdateJobRequest(req)
	if err != nil {
		js.log.Errorf("JobService.UpdateJob: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	_, err2 := js.jobDB.Update(ctx, job)
	if err2 != nil {
		js.log.Errorf("JobService.UpdateJob: %v", err2)
		return nil, grpcErrors.ErrorResponse(err2, err2.Error())
	}

	return &proto.UpdateJobResponse{Message: "Update success"}, nil
}

func (js *JobService) DeleteJob(ctx context.Context, req *proto.DeleteJobRequest) (*proto.DeleteJobResponse, error) {
	js.log.Infof("JobService.DeleteJob: grpc message arrived : %v", req)

	v := validator.DeleteJobRequestValidator{DeleteJobRequest: req}
	v_err := v.Validate()
	if v_err != nil {
		return nil, grpcErrors.ErrInvalidInput
	}

	err := js.jobDB.DeleteByID(ctx, req.GetId())
	if err != nil {
		js.log.Errorf("JobService.DeleteJob: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	return &proto.DeleteJobResponse{Message: "Delete success"}, nil
}
