package nats

import (
	"context"
	"encoding/json"
	"faas/internal/features/executions/domain/entity"
	"faas/internal/worker/domain/ports"
	"log"

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
	// Asegurarnos que el stream existe
	stream, err := c.js.StreamInfo("EXECUTIONS")
	if err != nil {
		log.Fatalf("Error getting stream info: %v", err)
	}
	log.Printf("Found stream EXECUTIONS with %d messages", stream.State.Msgs)

	// Configurar consumer
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

			msg.Ack()
			log.Printf("Successfully processed execution %s", execution.ID)
		},
		//nats.Durable(WORKERS_QUEUE),
		nats.ManualAck(),
		//nats.DeliverAll(),
		//nats.AckWait(time.Minute),
		//nats.MaxDeliver(3), // Reintentar hasta 3 veces
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
