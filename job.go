package gotask

import (
	"context"
	"fmt"

	"go.uber.org/atomic"
)

type Job struct {
	ctx      context.Context
	cancel   context.CancelFunc
	r        Runner
	doneHook func(job *Job)
	stopped  *atomic.Bool
	err      error
}

func newJob(hook func(job *Job)) *Job {
	return &Job{
		stopped:  atomic.NewBool(false),
		err:      nil,
		doneHook: hook,
	}
}

func (j *Job) done(err error) {
	j.err = err
	j.stopped.Store(true)
	if j.doneHook != nil {
		j.doneHook(j)
	}
}

func (j *Job) reset() {
	if j.cancel != nil {
		j.cancel()
	}
	j.ctx, j.cancel = context.WithCancel(context.TODO())
	j.stopped.Store(false)
	j.err = nil
}

func (j *Job) IsRunning() bool {
	return j.stopped.Load()
}

func (j *Job) Stop() {
	j.stopped.Store(true)
	if j.cancel != nil {
		j.cancel()
		j.cancel = nil
	}
	j.err = nil
}

func (j *Job) Err() error {
	return j.err
}

func (j *Job) Runner() Runner {
	if j.r != nil {
		return j.r
	}
	return nil
}

func (j *Job) String() string {
	key := ""
	if j.r != nil {
		key = fmt.Sprintf("%v", j.r.Key())
	}

	return "job:" + key
}

func runJob(job *Job, r Runner) (*Job, error) {
	job.reset()
	if r == nil {
		job.done(ErrNoRunners)
		return job, ErrNoRunners
	}
	job.r = r
	go func() {
		err := job.r.Run(job.ctx, job)
		job.done(err)
	}()
	return job, nil
}

func stopJob(t *task, job *Job, hook func(j *Job)) {
	t.moveIdleJob(job)
	if hook != nil {
		hook(job)
	}
	job.Stop()
}
