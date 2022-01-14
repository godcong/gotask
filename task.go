package gotask

import (
	"container/list"
	"sync"
)

type task struct {
	max  int
	lock *sync.RWMutex
	idle *sync.Pool
	jobs map[interface{}]*list.Element
	ll   *list.List
}

func (t *task) Start() error {
	return nil
}

type Task interface {
	Start() error
}

func Load(max int) Task {
	t := &task{
		max:  max,
		lock: &sync.RWMutex{},
		idle: &sync.Pool{},
		ll:   list.New(),
		jobs: make(map[interface{}]*list.Element),
		//doneHook: doneHook,
	}

	return t
}
