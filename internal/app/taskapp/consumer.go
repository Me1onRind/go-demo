package taskapp

import (
	"context"
	"net/http"

	"github.com/Me1onRind/go-demo/internal/infrastructure/async"
	"github.com/Shopify/sarama"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	JobManager *async.JobManager
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			defer session.MarkMessage(message, "")
			// handle
			var jobName string
			metadata := http.Header{}
			for _, header := range message.Headers {
				key := string(header.Key)
				value := string(header.Value)
				switch key {
				case async.KafkaJobNameKey:
					jobName = value
				default:
					metadata.Add(key, value)
				}
			}
			ctx := context.Background()
			job, _ := consumer.JobManager.GetJob(jobName)
			if job == nil {
				continue
			}
			_ = job.Handle(ctx, message.Value, metadata)
		case <-session.Context().Done():
			return nil
		}
	}
}
