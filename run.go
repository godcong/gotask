package gotask

import (
	"context"

	"github.com/google/uuid"
)

type Runner interface {
	Run(ctx context.Context, job *Job) error
	Key() interface{}
}

type KeyUUID struct{}

func (KeyUUID) Key() interface{} {
	return uuid.Must(uuid.NewRandom())
}
