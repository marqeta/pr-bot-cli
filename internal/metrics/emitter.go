package metrics

import (
	"context"
	"time"

	pmetrics "github.com/marqeta/pr-bot/metrics"
	httpmetrics "github.com/slok/go-http-metrics/metrics"
)

type NoopEmitter struct{}

func NewEmitter() pmetrics.Emitter {
	return &NoopEmitter{}
}

func (NoopEmitter) EmitDist(_ context.Context, _ string, _ float64, _ []string) {
}
func (NoopEmitter) EmitGauge(_ context.Context, _ string, _ float64, _ []string) {
}
func (NoopEmitter) ObserveHTTPRequestDuration(_ context.Context, _ httpmetrics.HTTPReqProperties, _ time.Duration) {
}
func (NoopEmitter) ObserveHTTPResponseSize(_ context.Context, _ httpmetrics.HTTPReqProperties, _ int64) {
}
func (NoopEmitter) AddInflightRequests(_ context.Context, _ httpmetrics.HTTPProperties, _ int) {
}
func (NoopEmitter) Close() {
}

var _ pmetrics.Emitter = &NoopEmitter{}
