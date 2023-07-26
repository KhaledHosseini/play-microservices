package handler

import (
	"net/http"
	"strconv"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/job/grpc"
	grpcutils "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/grpc"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	log logger.Logger
	*grpc.JobGRPCClient
}

func NewJobHandler(log logger.Logger, cfg *config.Config) *JobHandler {
	grpcClient := grpc.NewJobGRPCCLient(log, cfg) // Initialize the embedded type
	return &JobHandler{log: log, JobGRPCClient: grpcClient}
}

// @Summary Creates and schedule a job
// @Description Creates and schedule a job
// @Tags job
// @Accept json
// @Produce json
// @Param   input      body    models.CreateJobRequest     true        "input"
// @Success 200 {object} models.CreateJobResponse
// @Router /job/create [post]
func (jh *JobHandler) CreateJob(c *gin.Context) {
	jh.log.Info("JobHandler.CreateJob: Entered")
	var createJobRequest models.CreateJobRequest
	if errValid := c.ShouldBindJSON(&createJobRequest); errValid != nil {
		jh.log.Errorf("JobHandler.CreateJob: invalid CreateJobRequest")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "error": errValid.Error()})
		return
	}

	res, err := jh.GRPC_CreateJob(c.Request.Context(), createJobRequest.ToProto())
	if err != nil {
		jh.log.Errorf("JobHandler.CreateJob: error creation job: %v",err.Error())
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.CreateJobResponseFromProto(res))
}

// @Summary Updates and reschedule a job
// @Description Updates and reschedule a job
// @Tags job
// @Accept json
// @Produce json
// @Param   input      body    models.UpdateJobRequest     true        "input"
// @Success 200 {object} models.UpdateJobResponse
// @Router /job/update [post]
func (jh *JobHandler) UpdateJob(c *gin.Context) {
	jh.log.Info("JobHandler.UpdateJob: Entered")
	var updateJobRequest models.UpdateJobRequest
	if errValid := c.ShouldBindJSON(&updateJobRequest); errValid != nil {
		jh.log.Error("JobHandler.UpdateJob: invalid input updateJobRequest")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "error": errValid.Error()})
		return
	}

	res, err := jh.GRPC_UpdateJob(c.Request.Context(), updateJobRequest.ToProto())
	if err != nil {
		jh.log.Errorf("JobHandler.UpdateJob: error updating job: %v",err.Error())
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UpdateJobResponseFromProto(res))
}

// @Summary Get a job by id
// @Description Get a job by id
// @Tags job
// @Produce json
// @Param   id      query    string     true        "some id"
// @Success 200 {object} models.Job
// @Router /job/get [get]
func (jh *JobHandler) GetJob(c *gin.Context) {
	jh.log.Info("JobHandler.GetJob: Entered")
	id := c.Query("id")
	if id == "" {
		jh.log.Error("JobHandler.GetJob: invalid input. id not provided in query")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "error": "id not found"})
		return
	}

	res, err := jh.GRPC_GetJob(c.Request.Context(), &proto.GetJobRequest{Id: id})
	if err != nil {
		jh.log.Errorf("JobHandler.GetJob: error getting job: %v",err.Error())
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.JobFromProto(res))
}

// @Summary list jobs by type
// @Description list jobs by type
// @Tags job
// @Produce json
// @Param   page      query    int     true        "Page"
// @Param   size      query    int     true        "Size"
// @Success 200 {object} models.ListJobsResponse
// @Router /job/list [get]
func (jh *JobHandler) ListJobs(c *gin.Context) {
	jh.log.Info("JobHandler.ListJobs: Entered")
	page, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if err != nil {
		jh.log.Error("JobHandler.ListJobs: invalid input. page not provided in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "page parameter not provided"})
		return
	}

	size, err := strconv.ParseInt(c.Query("size"), 10, 32)
	if err != nil {
		jh.log.Error("JobHandler.ListJobs: invalid input. size not provided in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "size parameter not provided"})
		return
	}

	res, err := jh.GRPC_ListJobs(c.Request.Context(), &proto.ListJobsRequest{Page: page, Size: size})
	if err != nil {
		jh.log.Errorf("JobHandler.ListJobs: error listing jobs: %v",err.Error())
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ListJobsResponseFromProto(res))
}

// @Summary Delete a job by id
// @Description Delete a job by id
// @Tags job
// @Produce json
// @Param   id      query    string     true        "some id"
// @Success 200 {object} models.DeleteJobResponse
// @Router /job/delete [post]
func (jh *JobHandler) DeleteJob(c *gin.Context) {
	jh.log.Info("JobHandler.DeleteJob: Entered")
	id := c.Query("id")
	if id == "" {
		jh.log.Error("JobHandler.DeleteJob: invalid input. id is not provided in the query")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "error": "id not provided"})
		return
	}

	res, err := jh.GRPC_DeleteJob(c.Request.Context(), &proto.DeleteJobRequest{Id: id})
	if err != nil {
		jh.log.Errorf("JobHandler.DeleteJob: error deleting job: %v",err.Error())
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.DeleteJobResponseFromProto(res))
}
