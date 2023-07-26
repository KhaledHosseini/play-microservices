package handler

import (
	"net/http"
	"strconv"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models"
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/internal/models/report/grpc"
	grpcutils "github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/grpc"
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

// @Summary Get the list of reports
// @Description retrieve the reports
// @Tags report
// @Produce json
// @Param   page      query    int     true        "Page"
// @Param   size      query    int     true        "Size"
// @Success 200 {array} models.ListReportResponse
// @Router /report/list [get]
func (rh *ReportHandler) ListReports(c *gin.Context) {
	rh.log.Infof("Request arrived: list reports: %v", c.Request)

	page, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if err != nil {
		rh.log.Error("ReportHandler.ListReports: invalid input. page is not provided in the query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	size, err := strconv.ParseInt(c.Query("size"), 10, 32)
	if err != nil {
		rh.log.Error("ReportHandler.ListReports: invalid input. size is not provided in the query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
		return
	}

	res, err := rh.GRPC_ListReports(c.Request.Context(), &proto.ListReportsRequest{Page: page, Size: size})
	if err != nil {
		rh.log.Errorf("ReportHandler.ListReports: error listing reports: %v",err.Error())
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ListReportResponseFromProto(res))
}
