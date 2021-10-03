package demo_task

import (
	"github.com/Me1onRind/go-demo/internal/core/client/asynq_client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Send(t *testing.T) {
	asynq_client.InitAsynqClient("127.0.0.1:6379")
	err := SendDemoTask(&DemoTask{
		ID:   100,
		Name: "test",
	})
	assert.Empty(t, err)
}
