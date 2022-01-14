package gotask

type Runner interface {
	Run(job *Job) error
}

type Job struct {
	r Runner
}
