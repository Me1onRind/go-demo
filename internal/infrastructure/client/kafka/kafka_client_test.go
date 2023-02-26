package kafka

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SendJsonMessage(t *testing.T) {
	t.Run("sendDiffKey", func(t *testing.T) {
		cache := map[int32]struct{}{}
		client := NewMockKafkaClient(t)
		for i := 0; i < 100; i++ {
			client.MockProducer.ExpectSendMessageAndSucceed()
		}
		for i := 0; i < 100; i++ {
			partition, _, err := client.SendMessage(context.Background(), "test_topic", []byte("test"), PartitionKey(fmt.Sprintf("key_%d", i)))
			assert.Empty(t, err)
			cache[partition] = struct{}{}
		}
		assert.Equal(t, 4, len(cache))

	})
	t.Run("sendSameKeySuccess", func(t *testing.T) {
		client := NewMockKafkaClient(t)
		cache := map[int32]struct{}{}
		for i := 0; i < 20; i++ {
			client.MockProducer.ExpectSendMessageAndSucceed()
			partition, _, err := client.SendMessage(context.Background(), "test_topic", []byte("test2"), PartitionKey("key"))
			cache[partition] = struct{}{}
			assert.Empty(t, err)
		}
		assert.Equal(t, 1, len(cache))
	})
}
