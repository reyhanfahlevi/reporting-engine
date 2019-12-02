package reporting

import (
	"context"
	"fmt"

	es_report "github.com/tokopedia/reporting-engine/internal/model/es-report"
	"github.com/uniplaces/carbon"
)

type Reporting struct {
	esReportSvc esReportService
}

type esReportService interface {
	StoreReport(ctx context.Context, param es_report.ParamReporting) error
	GetReports(ctx context.Context, param es_report.ParamGetReports) ([]map[string]interface{}, int64, error)
}

// New will instantiate reporting package struct
func New(esReport esReportService) *Reporting {
	return &Reporting{
		esReportSvc: esReport,
	}
}

const (
	defaultMapping = `
	{
		"mappings": {
			"%s": {
				"properties": {
					"created_time": {
						"type": "date",
						"format": "yyyy-MM-dd HH:mm:ss"
					}
				}
			}
		}
	}`
)

// ParamSaveReport param for save report
type ParamSaveReport struct {
	ServiceName string
	DataType    string
	Data        map[string]interface{}
}

// SaveReport save report to service layer
func (r *Reporting) SaveReport(ctx context.Context, param ParamSaveReport) error {
	param.Data["created_time"] = carbon.Now().String()

	report := es_report.ParamReporting{
		Index:   param.ServiceName,
		Type:    param.DataType,
		Data:    param.Data,
		Mapping: fmt.Sprintf(defaultMapping, param.DataType),
	}

	return r.esReportSvc.StoreReport(ctx, report)
}

// ParamGetReports struct
type ParamGetReports struct {
	ServiceName string                 `json:"service_name"`
	DataType    string                 `json:"data_type"`
	Filter      map[string]interface{} `json:"filter"`
	RangeFilter map[string]RangeFilter `json:"range_filter"`
	Page        int                    `json:"page"`
	Limit       int                    `json:"limit"`
}

// RangeFilter struct
type RangeFilter struct {
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

// GetReportsResponse struct
type GetReportsResponse struct {
	Data       []map[string]interface{} `json:"data"`
	Pagination struct {
		PrevPage  int `json:"prev_page"`
		NextPage  int `json:"next_page"`
		TotalPage int `json:"total_page"`
	} `json:"pagination"`
}

// GetReports get reports data
func (r *Reporting) GetReports(ctx context.Context, param ParamGetReports) (GetReportsResponse, error) {
	var (
		resp GetReportsResponse
	)

	if param.Page <= 0 {
		param.Page = 1
	}

	if param.Limit <= 0 {
		param.Limit = 1
	}

	getParam := es_report.ParamGetReports{
		Index:  param.ServiceName,
		Type:   param.DataType,
		Filter: param.Filter,
	}

	for k, v := range param.RangeFilter {
		getParam.RangeFilter[k] = es_report.RangeFilter{
			From: v.From,
			To:   v.To,
		}
	}

	getParam.From = (param.Page - 1) * param.Limit
	getParam.Size = param.Limit

	reports, total, err := r.esReportSvc.GetReports(ctx, getParam)
	if err != nil {
		return resp, err
	}

	resp.Data = reports
	if param.Limit > 0 {
		totalPage := int(total) / param.Limit
		mod := int(total) % param.Limit
		if mod > 0 {
			totalPage++
		}
		resp.Pagination.TotalPage = totalPage
	}

	resp.Pagination.NextPage = param.Page + 1
	resp.Pagination.PrevPage = param.Page - 1

	if resp.Pagination.PrevPage <= 0 {
		resp.Pagination.PrevPage = 1
	}

	if resp.Pagination.NextPage > resp.Pagination.TotalPage {
		resp.Pagination.NextPage = resp.Pagination.TotalPage
	}

	return resp, nil
}
