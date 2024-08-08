package reportManager

import (
	"context"
	"github.com/marqeta/pr-bot/opa/input"

	"github.com/marqeta/pr-bot/opa/evaluation"
)

type ReportBuilder struct{}

func (r ReportBuilder) AddModuleResult(module string, result evaluation.Result) {
	return
}

func (r ReportBuilder) SetInput(input *input.Model) {
	return
}

func (r ReportBuilder) SetOutcome(result evaluation.Result) {
	return
}

func (r ReportBuilder) GetReport() evaluation.Report {
	return evaluation.Report{}
}

type ReportManager struct {
}

func (rm *ReportManager) NewReportBuilder(ctx context.Context, pr, reqID, deliveryID string) evaluation.ReportBuilder {
	return ReportBuilder{}
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
