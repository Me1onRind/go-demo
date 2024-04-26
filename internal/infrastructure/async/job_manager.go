package async

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/Me1onRind/go-demo/internal/global/gerror"
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
	jobs map[string]JobWorker // key: jobName
}

type MsgEntity struct {
	JobName string `json:"job_name"`
	Content any    `json:"content"`
}

func NewJobManager() *JobManager {
	return &JobManager{
		jobs: make(map[string]JobWorker, 0),
	}
}

func (j *JobManager) RegisterJob(job JobWorker) error {
	if _, ok := j.jobs[job.Name()]; ok {
		return fmt.Errorf("Register job failed, duplicate name:[%s]", job.Name())
	}
	j.jobs[job.Name()] = job
	return nil
}

func (j *JobManager) GetJob(name string) (JobWorker, error) {
	value, ok := j.jobs[name]
	if !ok {
		return nil, fmt.Errorf("%w, name:[%s]", ErrJobNotFound, name)
	}
	return value, nil
}

func (j *JobManager) Send(ctx context.Context, name string, protocol any, opts ...Option) error {
	job, err := j.GetJob(name)
	if err != nil {
		return err
	}
	if err := job.CheckSendProtocolType(ctx, protocol); err != nil {
		return err
	}
	sendParam := &SendParam{}
	for _, opt := range opts {
		opt(sendParam)
	}
	msgEntity := &MsgEntity{
		JobName: job.Name(),
		Content: protocol,
	}

	body, err := jsoniter.Marshal(msgEntity)
	if err != nil {
		return err
	}

	if err := job.Send(ctx, body, sendParam); err != nil {
		return err
	}

	return nil
}

func (j *JobManager) GetAllJobs() map[string]JobWorker {
	return j.jobs
}

type JobWorker interface {
	Send(ctx context.Context, msgEntity []byte, param *SendParam) error
	BackendType() JobBackendType
	Handle(ctx context.Context, body []byte, metadata http.Header) error

	// public method
	Name() string
	CheckSendProtocolType(context.Context, any) error
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

func (j *jobBase[T]) Name() string {
	return j.JobName
}

func (j *jobBase[T]) CheckSendProtocolType(ctx context.Context, protocol any) error {
	if _, ok := protocol.(*T); !ok {
		errMsg := fmt.Sprintf("Job %s send fail, protocol %s is not match register protocol %s", j.JobName, reflect.TypeOf(protocol), reflect.TypeOf(new(T)))
		logger.CtxErrorf(ctx, errMsg)
		return gerror.InvalidJobProtocolError.With(errMsg)
	}
	return nil
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
