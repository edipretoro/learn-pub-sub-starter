package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	c, _, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	deliveries, err := c.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		for d := range deliveries {
			var msg T
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Printf("could not unmarshal JSON message: %v", err)
				d.Nack(false, false)
				continue
			}
			handler(msg)
			d.Ack(false)
		}
	}()
	return nil
}
