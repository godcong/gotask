package gotask

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrRunnerIsNil = errors.New("runner is nil")
var ErrNoRunners = errors.New("no runners")

type Runner interface {
	Run(ctx context.Context, job *Job) error
	Key() interface{}
}

type KeyUUID struct {
	key interface{}
}

func (k *KeyUUID) Key() interface{} {
	if k.key == nil {
		k.key = uuid.Must(uuid.NewRandom())
	}
	return k.key
}
