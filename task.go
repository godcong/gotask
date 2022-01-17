package gotask

import (
	"container/list"
	"errors"
	"sync"
)

type task struct {
	max      int
	lock     *sync.RWMutex
	idle     *sync.Pool
	jobs     map[interface{}]*list.Element
	ll       *list.List
	doneHook func(job *Job)
}

var ErrTaskRunOverMax = errors.New("task run over max")
var ErrJobNotExists = errors.New("job not exists")

type Task interface {
	AddRunner(runner Runner) (*Job, error)
	StopJob(key interface{}) error
	Runs() int
}

func Load(max int, done func(j *Job)) Task {
	t := &task{
		max:  max,
		lock: &sync.RWMutex{},
		idle: &sync.Pool{},
		ll:   list.New(),
		jobs: make(map[interface{}]*list.Element),
		//doneHook: doneHook,
	}

	t.doneHook = func(job *Job) {
		stopJob(t, job, done)
	}
	return t
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

func (t *task) AddRunner(runner Runner) (*Job, error) {
	if t.IsFree() {
		return nil, ErrTaskRunOverMax
	}
	job, err := t.idleJob(runner)
	if err != nil {
		return nil, err
	}

	return runJob(job, runner)
}

func (t *task) StopJob(key interface{}) error {
	t.lock.RLock()
	ele, ok := t.jobs[key]
	t.lock.RUnlock()
	if !ok {
		return ErrJobNotExists
	}
	stopJob(t, ele.Value.(*Job), nil)
	return nil
}

func (t *task) idleJob(r Runner) (*Job, error) {
	ee := (*list.Element)(nil)
	ok := false
	t.lock.Lock()
	defer t.lock.Unlock()
	if ee, ok = t.jobs[r.Key()]; ok {
		t.ll.MoveToFront(ee)
		return ee.Value.(*Job), nil
	}
	vv := t.idle.Get()
	if vv == nil {
		vv = newJob(t.doneHook)
	}
	ee = t.ll.PushFront(vv)
	t.jobs[r.Key()] = ee
	return ee.Value.(*Job), nil
}

func (t *task) moveIdleJob(job *Job) {
	t.lock.Lock()
	if ee, ok := t.jobs[job.r.Key()]; ok {
		t.ll.Remove(ee)
		delete(t.jobs, job.r.Key())
		t.idle.Put(ee.Value.(*Job))
	}
	t.lock.Unlock()
}
