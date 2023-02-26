package async

import (
	"context"
	"fmt"
)

type JobManager struct {
	jobs map[string]Job // key: jobName
}

func (j *JobManager) RegisterJob(job Job) error {
	if _, ok := j.jobs[job.Name()]; ok {
		return fmt.Errorf("Register job failed, duplicate name:[%s]", job.Name())
	}
	j.jobs[job.Name()] = job
	return nil
}

type Job interface {
	Send(ctx context.Context, protocol any, opts ...Option) error
	Name() string
}

type SendParam struct {
	Key string
}

type Option func(*SendParam)

func WithKey(key string) Option {
	return func(s *SendParam) {
		s.Key = key
	}
}
