package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func Test_Send_Message(t *testing.T) {
	t.Skip()
	client, err := NewKafkaClient(configmd.KafkaConfig{
		Addr:            []string{},
		ProducerTimeout: time.Second * 2,
	})
	if !assert.Empty(t, err) {
		return
	}
	ctx := context.Background()
	err = client.SendMessage(ctx, &sarama.ProducerMessage{
		Value: sarama.StringEncoder("test data"),
	})
	assert.Empty(t, err)
}
