package handler

import (
	"net/http"
	"strconv"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/job/grpc"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type JobHandler struct {
	log logger.Logger
	*grpc.JobGRPCClient
}

func NewJobHandler(log logger.Logger, cfg *config.Config) *JobHandler {
	grpcClient := grpc.NewJobGRPCCLient(log, cfg) // Initialize the embedded type
	return &JobHandler{log: log, JobGRPCClient: grpcClient}
}

// @BasePath /api/v1/job
//
// @Summary Creates and schedule a job
// @Description Creates and schedule a job
// @Tags job
// @Accept json
// @Produce json
// @Success 200 {string} Job
// @Router /Create [post]
func (jh *JobHandler) CreateJob(c *gin.Context) {
	var createJobRequest models.CreateJobRequest
	if errValid := c.ShouldBindWith(&createJobRequest, binding.Form); errValid != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": errValid})
		return
	}

	res, err := jh.GRPC_CreateJob(c, createJobRequest.ToProto())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @BasePath /api/v1/job
//
// @Summary Get a job
// @Description Get a job
// @Tags job
// @Param
// @Produce json
// @Success 200 {string} Job
// @Router /Get [get]
func (jh *JobHandler) GetJob(c *gin.Context) {
	var getJobRequest models.GetJobRequest
	if errValid := c.ShouldBindWith(&getJobRequest, binding.Form); errValid != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": errValid})
		return
	}

	res, err := jh.GRPC_GetJob(c, getJobRequest.ToProto())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (jh *JobHandler) ListJobs(c *gin.Context) {
	jh.log.Infof("Request arrived: list jobs: %v", c.Request)
	page, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if err != nil {
		// Handle missing or invalid page parameter
		c.JSON(400, gin.H{"error": "Invalid page parameter"})
		return
	}

	size, err := strconv.ParseInt(c.Query("size"), 10, 32)
	if err != nil {
		// Handle missing or invalid page parameter
		c.JSON(400, gin.H{"error": "Invalid size parameter"})
		return
	}

	res, err := jh.GRPC_ListJobs(c, &proto.ListJobsRequest{Page: page, Size: size})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (jh *JobHandler) UpdateJob(c *gin.Context) {
	var updateJobRequest models.UpdateJobRequest
	if errValid := c.ShouldBindWith(&updateJobRequest, binding.Form); errValid != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": errValid})
		return
	}

	res, err := jh.GRPC_UpdateJob(c, updateJobRequest.ToProto())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (jh *JobHandler) DeleteJob(c *gin.Context) {
	var deleteJobRequest models.DeleteJobRequest
	if errValid := c.ShouldBindWith(&deleteJobRequest, binding.Form); errValid != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": errValid})
		return
	}

	res, err := jh.GRPC_DeleteJob(c, deleteJobRequest.ToProto())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err})
		return
	}

	c.JSON(http.StatusOK, res)
}
