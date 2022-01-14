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

func (t *task) IsFree() bool {
	return t.Runs() >= t.max
}

func (t *task) IsRunning(key interface{}) bool {
	t.lock.RLock()
	_, ok := t.jobs[key]
	t.lock.RUnlock()
	return ok
}

func (t *task) Runs() (i int) {
	t.lock.RLock()
	i = t.ll.Len()
	t.lock.RUnlock()
	return
}

func (t *task) AddRunner(runner Runner) error {
	//if t.IsFree() {
	//	return ErrTaskRunOverMax
	//}
	//
	//job, err := t.idleJob(state)
	//if err != nil {
	//	return err
	//}
	//if err := runJob(t.api, job, state); err != nil {
	//	t.moveIdleJob(job)
	//}
	//return nil
	return nil
}

func (t *task) Start() error {
	return nil
}

type Task interface {
	Start() error
	AddRunner(runner Runner) error
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
