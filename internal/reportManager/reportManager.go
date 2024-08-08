package reportManager

import (
	"context"

	"github.com/marqeta/pr-bot/opa/evaluation"
)

type liteReportBuilder struct{}

type ReportManager struct {
}

func (rm *ReportManager) NewReportBuilder(ctx context.Context, pr, reqID, deliveryID string) evaluation.ReportBuilder {
	return nil
}

func (rm *ReportManager) GetReport(ctx context.Context, pr, deliveryID string) (*evaluation.Report, error) {
	return nil, nil
}

func (rm *ReportManager) StoreReport(ctx context.Context, builder evaluation.ReportBuilder) error {
	return nil
}

func (rm *ReportManager) ListReports(ctx context.Context, pr string) ([]evaluation.ReportMetadata, error) {
	return nil, nil
}
