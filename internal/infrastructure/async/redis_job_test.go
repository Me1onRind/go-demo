package async

import (
	"context"
	"testing"

	"github.com/Me1onRind/go-demo/internal/infrastructure/unittest"
	"github.com/stretchr/testify/assert"
)

func TestSendRedisJob(t *testing.T) {
	unittest.SetRedisJobConfig("demo", "redis_demo", "queue_demo")

	type Msg struct {
		Field string
	}

	t.Run("sendSuccess", func(t *testing.T) {
		redisServer := unittest.GetMockRedis(t, "redis_demo")
		defer redisServer.Close()

		jobWorker := NewRedisJob[Msg]("demo", nil)
		jm := NewJobManager()
		jm.RegisterJob(jobWorker)

		err := jm.Send(context.Background(), "demo", &Msg{Field: "field"})
		assert.Empty(t, err)
		redisServer.CheckList(t, "queue_demo", `{"job_name":"demo","content":{"Field":"field"}}`)
	})
}
