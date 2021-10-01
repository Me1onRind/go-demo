package asynq_client

import (
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
)

func Test_SendTask(t *testing.T) {
	InitAsynqClient("127.0.0.1:6379")
	defer AsynqClient.Close()
	_, err := AsynqClient.Enqueue(asynq.NewTask("test", []byte("paylod")))
	assert.Empty(t, err)
}
