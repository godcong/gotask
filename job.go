package gotask

import (
	"context"
)

type Runner interface {
	Run(ctx context.Context, job *Job) error
}

type Job struct {
	r Runner
}
