package async

import (
	"context"
	"errors"
	"testing"

	"github.com/Me1onRind/go-demo/internal/infrastructure/unittest"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func Test_KafkaJob_Send(t *testing.T) {
	type Msg struct {
		Field string
	}

	kafkaClient := unittest.GetMockKafkaClient(t, "kafka_demo")
	t.Log(kafkaClient)
	unittest.SetKafkaJobConfig("demo", "kafka_demo", "topci_demo")
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
}
