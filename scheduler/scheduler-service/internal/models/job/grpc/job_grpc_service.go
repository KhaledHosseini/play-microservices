package grpc

import (
	context "context"
	"fmt"
	"time"

	models "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models"
	validator "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models/job/validation"
	jobFucntion "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/job_function"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/logger"
	proto "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/proto"

	scheduler "github.com/reugn/go-quartz/quartz"
	"github.com/segmentio/fasthash/fnv1a"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// load all pending jobs from database to the scheduler again

}

func (js *JobService) scheduleJob(ctx context.Context, key int, jobID string, scheduleTime time.Time) {

	_ = js.jobsScheduler.DeleteJob(key)

	functionJob := jobFucntion.NewFunctionJobWithKey(key, func(parentCtx context.Context) (int, error) {
		ctx, cancel := context.WithTimeout(parentCtx, 3*time.Second)
		defer cancel()

		js.log.Info("Job is triggering.")
		jb, err := js.jobDB.GetByID(ctx, jobID)
		if err != nil {
			js.log.Errorf("JobService.JobFunction.GetByID: %v", err.Error())
			return 0, status.Errorf(codes.NotFound, "Job not found in the database.")
		}
		//set the state of the job to running
		jb.Status = int32(proto.JobStatus_JOB_STATUS_RUNNING)
		_, err2 := js.jobDB.Update(ctx, jb)
		if err2 != nil {
			js.log.Errorf("JobService.JobFunction.Update: %v", err)
			return 0, status.Errorf(codes.NotFound, "cannot update job status")
		}

		js.jobsProducer.PublishRun(ctx, jb)
		js.log.Info("Job triggered and run has published.")
		return 0, nil
	})
	currentTime := time.Now()
	duration := scheduleTime.Sub(currentTime)
	if duration <= 0 {
		duration = 10 * time.Second
	}
	js.log.Infof("trigger time is %s ", scheduleTime)
	js.log.Infof("triggers in %s seconds", duration)
	js.jobsScheduler.ScheduleJob(ctx, functionJob, scheduler.NewRunOnceTrigger(duration))
}
func (js *JobService) CreateJob(ctx context.Context, req *proto.CreateJobRequest) (*proto.CreateJobResponse, error) {
	js.log.Infof("JobService.CreateJob: grpc message arrived : %v", req)
	v := validator.CreateJobRequestValidator{CreateJobRequest: req}
	v_err := v.Validate()
	if v_err != nil {
		js.log.Errorf("JobService.CreateJob: input validation failed for request : %v with error %v", req, v_err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request body.")
	}

	job, err0 := models.JobFromProto_CreateJobRequest(req)
	if err0 != nil {
		js.log.Errorf("JobService.CreateJob: cannot create proto request")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request body.")
	}
	js.log.Infof("JobService.CreateJob: proto converted to db model : %v", job)
	job.Status = int32(proto.JobStatus_JOB_STATUS_SCHEDULED)
	jobFingerPrint := fmt.Sprintf("%s:%s:%s:%p", job.Name, job.Description, job.JobData, &job.ScheduleTime)
	job.ScheduledKey = int(fnv1a.HashString64(jobFingerPrint)) //We assume 64 bit systems!
	//store the job in the database
	js.log.Infof("JobService.CreateJob: Creating job in the database")
	created, err := js.jobDB.Create(ctx, job)
	if err != nil {
		js.log.Errorf("JobService.CreateJob: %v", err.Error())
		return nil, status.Errorf(codes.Internal, "Cannot save job. db error")
	}

	jobID := created.Id.Hex()
	js.scheduleJob(ctx, job.ScheduledKey, jobID, job.ScheduleTime)
	js.jobsProducer.PublishCreate(ctx, created)

	return &proto.CreateJobResponse{Id: jobID}, nil
}

func (js *JobService) GetJob(ctx context.Context, req *proto.GetJobRequest) (*proto.Job, error) {
	js.log.Infof("JobService.GetJob: grpc message arrived : %v", req)

	v := validator.GetJobRequestValidator{GetJobRequest: req}
	v_err := v.Validate()
	if v_err != nil {
		js.log.Errorf("JobService.GetJob: input validation failed for request : %v with error %v", req, v_err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request body.")
	}

	job, err := js.jobDB.GetByID(ctx, req.GetId())
	if err != nil {
		js.log.Errorf("JobService.GetJob: %v", err.Error())
		return nil, status.Errorf(codes.Internal, "cannot retrieve job from db.")
	}

	return job.ToProto(), nil
}

func (js *JobService) ListJobs(ctx context.Context, req *proto.ListJobsRequest) (*proto.ListJobsResponse, error) {
	js.log.Infof("JobService.ListJobs: grpc message arrived : %v", req)

	v := validator.ListJobsRequestValidator{ListJobsRequest: req}
	v_err := v.Validate()
	if v_err != nil {
		js.log.Errorf("JobService.ListJobs: input validation failed for request : %v with error %v", req, v_err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request body.")
	}

	jobs, err := js.jobDB.ListALL(ctx, req.Size, req.Page)
	if err != nil {
		js.log.Errorf("JobService.ListJobs: %v", err)
		return nil, status.Errorf(codes.Internal, "Cannot get list from db.")
	}

	js.log.Infof("JobService.ListJobs: jobs returned: %v", jobs)

	return jobs.ToProto(), nil
}

func (js *JobService) UpdateJob(ctx context.Context, req *proto.UpdateJobRequest) (*proto.UpdateJobResponse, error) {
	js.log.Infof("JobService.UpdateJob: grpc message arrived : %v", req)

	v := validator.UpdateJobRequestValidator{UpdateJobRequest: req}
	v_err := v.Validate()
	if v_err != nil {
		js.log.Errorf("JobService.UpdateJob: input validation failed for request : %v with error %v", req, v_err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request body.")
	}

	job, err := models.JobFromProto_UpdateJobRequest(req)
	if err != nil {
		js.log.Errorf("JobService.UpdateJob: cannot create proto with error: %v", err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request body.")
	}

	_, err2 := js.jobDB.Update(ctx, job)
	if err2 != nil {
		js.log.Errorf("JobService.UpdateJob: %v", err2)
		return nil, status.Errorf(codes.Internal, "cannot update job")
	}

	jobID := job.Id.Hex()
	js.scheduleJob(ctx, job.ScheduledKey, jobID, job.ScheduleTime)

	js.jobsProducer.PublishUpdate(ctx, job)

	return &proto.UpdateJobResponse{Message: "Update success"}, nil
}

func (js *JobService) DeleteJob(ctx context.Context, req *proto.DeleteJobRequest) (*proto.DeleteJobResponse, error) {
	js.log.Infof("JobService.DeleteJob: grpc message arrived : %v", req)

	v := validator.DeleteJobRequestValidator{DeleteJobRequest: req}
	v_err := v.Validate()
	if v_err != nil {
		js.log.Errorf("JobService.DeleteJob: input validation failed for request : %v with error %v", req, v_err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request body.")
	}

	err := js.jobDB.DeleteByID(ctx, req.GetId())
	if err != nil {
		js.log.Errorf("JobService.DeleteJob: %v", err)
		return nil, status.Errorf(codes.Internal, "cannot delete job.")
	}

	return &proto.DeleteJobResponse{Message: "Delete success"}, nil
}
