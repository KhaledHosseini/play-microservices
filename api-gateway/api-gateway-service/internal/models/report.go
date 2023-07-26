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
	Topic        string
	CreatedTime time.Time
	ReportData  string
}

func ReportFromProto(p *proto.Report) *Report {
	return &Report{
		Id:           p.Id,
		Topic:        p.Topic,
		CreatedTime: p.CreatedTime.AsTime(),
		ReportData:  p.ReportData,
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
