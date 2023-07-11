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
	var createJobRequest models.CreateJobRequest
	if errValid := c.ShouldBindJSON(&createJobRequest); errValid != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": errValid.Error()})
		return
	}

	res, err := jh.GRPC_CreateJob(c, createJobRequest.ToProto())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.CreateJobResponseFromProto(res))
}

// @Summary Get a job by id
// @Description Get a job by id
// @Tags job
// @Produce json
// @Param   id      path    string     true        "some id"
// @Success 200 {object} models.Job
// @Router /job/{id} [get]
func (jh *JobHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": "id not found"})
		return
	}

	res, err := jh.GRPC_GetJob(c, &proto.GetJobRequest{Id: id})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
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

	c.JSON(http.StatusOK, models.ListJobsResponseFromProto(res))
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
	var updateJobRequest models.UpdateJobRequest
	if errValid := c.ShouldBindJSON(&updateJobRequest); errValid != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": errValid.Error()})
		return
	}

	res, err := jh.GRPC_UpdateJob(c, updateJobRequest.ToProto())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UpdateJobResponseFromProto(res))
}

// @Summary Delete a job by id
// @Description Delete a job by id
// @Tags job
// @Produce json
// @Param   id      path    string     true        "some id"
// @Success 200 {object} models.DeleteJobResponse
// @Router /job/{id} [delete]
func (jh *JobHandler) DeleteJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"status": false, "error": "id not provided"})
		return
	}

	res, err := jh.GRPC_DeleteJob(c, &proto.DeleteJobRequest{Id: id})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.DeleteJobRequestFromProto(res))
}
