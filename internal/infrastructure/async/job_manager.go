package async

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Me1onRind/go-demo/internal/infrastructure/goroutine"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/infrastructure/middleware"
	jsoniter "github.com/json-iterator/go"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	ErrJobNotFound = errors.New("Job is not found")
)

type JobManager struct {
	jobs map[string]Job // key: jobName
}

func NewJobManager() *JobManager {
	return &JobManager{
		jobs: make(map[string]Job, 0),
	}
}

func (j *JobManager) RegisterJob(job Job) error {
	if _, ok := j.jobs[job.Name()]; ok {
		return fmt.Errorf("Register job failed, duplicate name:[%s]", job.Name())
	}
	j.jobs[job.Name()] = job
	return nil
}

func (j *JobManager) GetJob(name string) (Job, error) {
	value, ok := j.jobs[name]
	if !ok {
		return nil, fmt.Errorf("%w, name:[%s]", ErrJobNotFound, name)
	}
	return value, nil
}

func (j *JobManager) GetAllJobs() map[string]Job {
	return j.jobs
}

type Job interface {
	Send(ctx context.Context, protocol any, opts ...Option) error
	Name() string
	BackendType() JobBackendType
	Handle(ctx context.Context, body []byte, metadata http.Header) error
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

type jobBase[T any] struct {
	JobName string
	Handler func(context.Context, *T) error
}

func (j *jobBase[T]) Handle(ctx context.Context, body []byte, metadata http.Header) error {
	// recover
	defer func() {
		if err := recover(); err != nil {
			goroutine.LogPanicStack(ctx, err)
		}
	}()
	// trace
	ctx, span := middleware.Tracer(ctx, j.JobName, opentracing.HTTPHeadersCarrier(metadata))
	defer span.Finish()
	// log
	start := time.Now()
	logger.CtxInfof(ctx, "%s|%s|body=%s", j.JobName, time.Since(start), body)

	var p T
	if err := jsoniter.Unmarshal(body, &p); err != nil {
		errMsg := fmt.Sprintf("Job:[%s] handle fail, body:[%+v] umarshal fail", j.JobName, p)
		logger.CtxErrorf(ctx, errMsg)
		return errors.New(errMsg)
	}

	return j.Handler(ctx, &p)
}
