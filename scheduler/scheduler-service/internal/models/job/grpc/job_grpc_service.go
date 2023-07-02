package grpc

import (
	context "context"
	"fmt"
	"math"
	"time"

	models "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models"
	grpcErrors "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/grpc_errors"
	jobFucntion "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/job_function"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/logger"
	pb "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/proto"

	scheduler "github.com/reugn/go-quartz/quartz"
	"github.com/segmentio/fasthash/fnv1a"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type JobService struct {
	log           logger.Logger
	jobDB         models.JobDB
	jobsProducer  models.JobsProducer
	jobsScheduler scheduler.Scheduler
	pb.UnimplementedJobServiceServer
}

func NewJobService(log logger.Logger, jobDB models.JobDB, jobsProducer models.JobsProducer, jobsScheduler scheduler.Scheduler) *JobService {
	return &JobService{log: log, jobDB: jobDB, jobsProducer: jobsProducer, jobsScheduler: jobsScheduler}
}

// exported method (Started with uppercase)
func (js *JobService) LoadScheduledJobs() {
	// load all pending jobs to the scheduler again
}

func (j *JobService) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {

	job := models.JobFromProto_CreateJobRequest(req)
	job.Status = int32(pb.JobStatus_SCHEDULED)
	jobFingerPrint := fmt.Sprintf("%s:%s:%s:%p", job.Name, job.Description, job.JobData, &job.ScheduleTime)
	job.ScheduledKey = int(fnv1a.HashString64(jobFingerPrint)) //We assume 64 bit systems!
	//cretae job in the database
	created, err := j.jobDB.Create(ctx, job)
	if err != nil {
		j.log.Errorf("JobService.CreateJob: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	jobID := created.JobID
	functionJob := jobFucntion.NewFunctionJobWithKey(job.ScheduledKey, func(_ context.Context) (int, error) {

		jb, err := j.jobDB.GetByID(ctx, jobID)
		if err != nil {
			j.log.Errorf("JobService.JobFunction.GetByID: %v", err)
			return 0, grpcErrors.ErrorResponse(err, err.Error())
		}
		//set the state of the job to running
		jb.Status = int32(pb.JobStatus_RUNNING)
		_, err2 := j.jobDB.Update(ctx, jb)
		if err2 != nil {
			j.log.Errorf("JobService.JobFunction.Update: %v", err)
			return 0, grpcErrors.ErrorResponse(err2, err2.Error())
		}

		j.jobsProducer.PublishRun(ctx, jb)

		return 0, nil
	})

	triggerTime := time.Duration(math.Max(0, time.Until(job.ScheduleTime).Seconds()))
	j.log.Infof("triggers in %s seconds", triggerTime)

	j.jobsScheduler.ScheduleJob(ctx, functionJob, scheduler.NewRunOnceTrigger(triggerTime))

	j.jobsProducer.PublishCreate(ctx, created)

	return &pb.CreateJobResponse{Id: created.JobID.String()}, nil
}

func (js *JobService) GetJob(context.Context, *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJob not implemented")
}

func (js *JobService) ListJobs(context.Context, *pb.ListJobsRequest) (*pb.ListJobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListJobs not implemented")
}

func (js *JobService) UpdateJob(context.Context, *pb.UpdateJobRequest) (*pb.UpdateJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateJob not implemented")
}

func (js *JobService) DeleteJob(context.Context, *pb.DeleteJobRequest) (*pb.DeleteJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteJob not implemented")
}
