package gotask

import (
	"context"

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

func runJob(job *Job, r Runner) error {
	job.reset()
	if r == nil {
		job.done(ErrNoRunners)
		return nil
	}
	job.r = r
	go func() {
		err := job.r.Run(job.ctx, job)
		job.done(err)
	}()
	return nil
}

func stopJob(t *task, job *Job, hook func(j *Job)) {
	t.moveIdleJob(job)
	if hook != nil {
		hook(job)
	}
	job.Stop()
}
