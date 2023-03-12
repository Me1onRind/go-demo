package async

import (
	"context"
	"errors"
	"testing"

	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/unittest"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func Test_KafkaJob_Send(t *testing.T) {
	type Msg struct {
		Field string
	}
	unittest.SetKafkaJobConfig("demo", "kafka_demo", "topic_demo")

	t.Run("sendSuccess", func(t *testing.T) {
		kafkaClient := unittest.GetMockKafkaClient(t, "kafka_demo")
		t.Log(kafkaClient)
		kafkaClient.MockProducer.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(msg *sarama.ProducerMessage) error {
			value, _ := msg.Value.Encode()
			if string(value) != `{"job_name":"demo","content":{"Field":"field"}}` {
				return errors.New("value not expect")
			}
			return nil
		})

		job := NewKafkaJob[Msg]("demo", nil)
		err := job.Send(context.Background(), &Msg{Field: "field"}, WithKey("key"))
		assert.Empty(t, err)
	})

	t.Run("sendWrongProtocol", func(t *testing.T) {
		type WrongMsg struct{}
		job := NewKafkaJob[Msg]("demo", nil)
		err := job.Send(context.Background(), &WrongMsg{}, WithKey("key"))
		assert.ErrorIs(t, err, gerror.SendJobError)
	})
}
