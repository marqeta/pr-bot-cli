package reportmanager

import (
	"context"

	"github.com/marqeta/pr-bot/opa/evaluation"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/rs/zerolog/log"
)

type ReportBuilder struct{}

func (r ReportBuilder) AddModuleResult(module string, result evaluation.Result) {
	log.Info().Str("module", module).Interface("result", result).Msg("Adding module result")
}

func (r ReportBuilder) SetInput(input *input.Model) {
	log.Info().Interface("input", input).Msg("Setting input")
}

func (r ReportBuilder) SetOutcome(result evaluation.Result) {
	log.Info().Interface("result", result).Msg("Setting outcome")
}

func (r ReportBuilder) GetReport() evaluation.Report {
	return evaluation.Report{}
}

type ReportManager struct {
}

func (rm *ReportManager) NewReportBuilder(_ context.Context, _, _, _ string) evaluation.ReportBuilder {
	return ReportBuilder{}
}

func (rm *ReportManager) GetReport(_ context.Context, _, _ string) (*evaluation.Report, error) {
	return nil, nil
}

func (rm *ReportManager) StoreReport(_ context.Context, _ evaluation.ReportBuilder) error {
	return nil
}

func (rm *ReportManager) ListReports(_ context.Context, _ string) ([]evaluation.ReportMetadata, error) {
	return nil, nil
}
