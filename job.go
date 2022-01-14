package gotask

import (
	"context"
)

type Runner interface {
	Run(ctx context.Context, job *Job) error
	Key() interface{}
}

type Job struct {
	r Runner
}
