package opa

import "context"

type Evaluator interface {
	Evaluate(ctx context.Context) error
}
