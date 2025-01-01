package nats

import (
	"context"
	"encoding/json"
	"faas/internal/features/executions/domain/entity"
	"faas/internal/worker/domain/ports"
	"log"
	"math"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	EXECUTIONS_SUBJECT = "executions.pending"
	WORKERS_QUEUE      = "execution-workers"
)

type NatsStreamConsumer struct {
	js           nats.JetStreamContext
	subscription *nats.Subscription
}

func NewStreamConsumer(js nats.JetStreamContext) ports.StreamConsumer {
	return &NatsStreamConsumer{js: js}
}

func (c *NatsStreamConsumer) Subscribe(handler func(ctx context.Context, execution *entity.Execution) error) ports.Worker {
	maxRetries := 5
	var stream *nats.StreamInfo
	var err error

	// Try to connect to stream with retry
	for i := 0; i < maxRetries; i++ {
		stream, err = c.js.StreamInfo("EXECUTIONS")
		if err != nil {
			waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
			log.Printf("Stream not available, retrying in %v seconds... (attempt %d/%d)",
				waitTime, i+1, maxRetries)
			time.Sleep(waitTime)
			continue
		}
		// Stream found
		log.Printf("Successfully connected to EXECUTIONS stream with %d messages", stream.State.Msgs)
		break
	}

	// If after all attempts it still fails
	if err != nil {
		log.Fatalf("Failed to connect to stream after %d attempts: %v", maxRetries, err)
	}

	// Configure consumer
	sub, err := c.js.QueueSubscribe(
		EXECUTIONS_SUBJECT,
		WORKERS_QUEUE,
		func(msg *nats.Msg) {
			log.Printf("Received message: %s", string(msg.Data))
			var execution entity.Execution
			if err := json.Unmarshal(msg.Data, &execution); err != nil {
				log.Printf("Error unmarshaling execution: %v", err)
				msg.Nak()
				return
			}

			if err := handler(context.Background(), &execution); err != nil {
				log.Printf("Error processing execution: %v", err)
				msg.Nak()
				return
			}

			if err := msg.Ack(); err != nil {
				log.Printf("Error acknowledging message: %v", err)
			}

			log.Printf("Successfully processed execution %s", execution.ID)
		},
		nats.ManualAck(),
		nats.AckWait(30*time.Minute),
		nats.MaxDeliver(1),
		nats.DeliverAll(),
	)

	if err != nil {
		log.Fatalf("Error subscribing to stream: %v", err)
	}

	log.Printf("Successfully subscribed to %s with queue %s", EXECUTIONS_SUBJECT, WORKERS_QUEUE)
	c.subscription = sub
	return c
}

func (c *NatsStreamConsumer) Stop() error {
	if c.subscription != nil {
		return c.subscription.Unsubscribe()
	}
	return nil
}
