package evaluation

import (
	"context"

	peval "github.com/marqeta/pr-bot/opa/evaluation"
)

type NoopManager struct{}

var _ peval.Manager = (*NoopManager)(nil)

func (n NoopManager) NewReportBuilder(_ context.Context, _, _, _ string) peval.ReportBuilder {
	return &NoopReportBuilder{}
}

func (n NoopManager) GetReport(_ context.Context, _, _ string) (*peval.Report, error) {
	return &peval.Report{}, nil
}

func (n NoopManager) StoreReport(_ context.Context, _ peval.ReportBuilder) error {
	return nil
}

func (n NoopManager) ListReports(_ context.Context, _ string) ([]peval.ReportMetadata, error) {
	return nil, nil
}
