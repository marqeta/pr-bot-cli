package evaluation

import (
	peval "github.com/marqeta/pr-bot/opa/evaluation"
	"github.com/marqeta/pr-bot/opa/input"
)

type NoopReportBuilder struct{}

var _ peval.ReportBuilder = (*NoopReportBuilder)(nil)

func (n NoopReportBuilder) AddModuleResult(_ string, _ peval.Result) {
	// do nothing
}

func (n NoopReportBuilder) SetInput(_ *input.Model) {
	// do nothing
}

func (n NoopReportBuilder) SetOutcome(_ peval.Result) {
	// do nothing
}

func (n NoopReportBuilder) GetReport() peval.Report {
	return peval.Report{}
}
