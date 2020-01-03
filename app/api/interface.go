package api

import (
	"context"

	"github.com/tokopedia/td-report-engine/internal/usecase/reporting"
)

// ReportingUseCase reporting usecase contract
type ReportingUseCase interface {
	SaveReport(ctx context.Context, param reporting.ParamSaveReport) error
	GetReports(ctx context.Context, param reporting.ParamGetReports) (reporting.GetReportsResponse, error)
}
