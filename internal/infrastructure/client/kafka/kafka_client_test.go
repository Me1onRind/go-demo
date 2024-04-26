package kafka

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	t.Run("sendDiffKey", func(t *testing.T) {
		cache := map[int32]struct{}{}
		client := NewMockKafkaClient(t)
		for i := 0; i < 100; i++ {
			client.MockProducer.ExpectSendMessageAndSucceed()
		}
		for i := 0; i < 100; i++ {
			partition, _, err := client.SendMessage(context.Background(), "test_topic", fmt.Sprintf("key_%d", i), []byte("test"))
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
			partition, _, err := client.SendMessage(context.Background(), "test_topic", "key", []byte("test2"))
			cache[partition] = struct{}{}
			assert.Empty(t, err)
		}
		assert.Equal(t, 1, len(cache))
	})

	t.Run("sendRightValue", func(t *testing.T) {
		client := NewMockKafkaClient(t)
		client.MockProducer.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(pm *sarama.ProducerMessage) error {
			val, err := pm.Value.Encode()
			if err != nil {
				return err
			}
			if string(val) != "test2" {
				return errors.New("Send wrong message body")
			}

			if len(pm.Headers) != 1 {
				return errors.New("Send wrong headers length")
			}

			if string(pm.Headers[0].Key) != "request_id" || string(pm.Headers[0].Value) != "test123" {
				return errors.New("Send wrong header key or value")
			}

			return nil
		})
		_, _, err := client.SendMessage(context.Background(), "test_topic", "key", []byte("test2"), WithHeaers(map[string]string{
			"request_id": "test123",
		}))
		assert.Empty(t, err)
	})

	t.Run("sendMessageFail", func(t *testing.T) {
		client := NewMockKafkaClient(t)
		client.MockProducer.ExpectSendMessageAndFail(errors.New("mock failed"))

		_, _, err := client.SendMessage(context.Background(), "test_topic", "key", []byte("test2"), WithHeaers(map[string]string{
			"request_id": "test123",
		}))
		assert.ErrorIs(t, err, gerror.SendKafkaError)
	})

}
