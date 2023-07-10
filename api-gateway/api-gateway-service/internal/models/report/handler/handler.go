package handler

import (
	"net/http"
	"strconv"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/report/grpc"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	log logger.Logger
	*grpc.ReportGRPCClient
}

func NewReportHandler(log logger.Logger, cfg *config.Config) *ReportHandler {
	grpcClient := grpc.NewReportGRPCClient(log, cfg) // Initialize the embedded type
	return &ReportHandler{log: log, ReportGRPCClient: grpcClient}
}

// @BasePath /api/v1/job
//
// @Summary Creates and schedule a job
// @Description Creates and schedule a job
// @Tags job
// @Accept json
// @Produce json
// @Success 200 {string} Job
// @Router /Create [get]
func (rh *ReportHandler) ListReports(c *gin.Context) {
	rh.log.Infof("Request arrived: list reports: %v", c.Request)
	filter, err := strconv.ParseInt(c.Query("filter"), 10, 32)
	if err != nil {
		// Handle missing or invalid page parameter
		c.JSON(400, gin.H{"error": "Invalid filter parameter"})
		return
	}

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

	res, err := rh.GRPC_ListReports(c, &proto.ListReportsRequest{Filter: proto.ReportType(filter), Page: page, Size: size})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
