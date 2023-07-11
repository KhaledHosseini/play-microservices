package models

import (
	"time"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/thoas/go-funk"
)

type ListReportsRequest struct {
	Filter int32
	Page   int64
	Size   int64
}

type Report struct {
	Id           string
	Type         int32
	Topic        string
	Created_time time.Time
	Report_data  string
}

func ReportFromProto(p *proto.Report) *Report {
	return &Report{
		Id:           p.Id,
		Type:         int32(p.Type),
		Topic:        p.Topic,
		Created_time: p.CreatedTime.AsTime(),
		Report_data:  p.ReportData,
	}
}

type ListReportResponse struct {
	TotalCount int64
	TotalPages int64
	Page       int64
	Size       int64
	HasMore    bool
	Reports    []Report
}

func ListReportResponseFromProto(p *proto.ListReportResponse) *ListReportResponse {
	reports := funk.Map(p.Reports, func(x *proto.Report) Report {
		return *ReportFromProto(x)
	}).([]Report)

	return &ListReportResponse{
		TotalCount: p.TotalCount,
		TotalPages: p.TotalPages,
		Page:       p.Page,
		Size:       p.Size,
		HasMore:    p.HasMore,
		Reports:    reports,
	}
}
