package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type rabbitMQService struct{
	conn *amqp091.Connection
	ch *amqp091.Channel
	logger *zerolog.Logger
}

func NewRabbitMQService(amqpURL string, logger *zerolog.Logger) (RabbitMQService, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		logger.Error().Err(err).Msgf("Failed to connect to RabbitMQ")
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		logger.Error().Err(err).Msgf("Failed to create channel")
		return nil, err
	}

	return &rabbitMQService{
		conn: conn,
		ch: ch,
		logger: logger,
	}, nil

}

func (r *rabbitMQService) Publish(ctx context.Context, queue string, message any) error {
	_, err := r.ch.QueueDeclare(
		queue, // name
		true,    // durability
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		amqp091.Table{
		  amqp091.QueueTypeArg: amqp091.QueueTypeQuorum,
		},
	  )
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to declare queue %s", queue)
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to marshal message")
		return err
	}
err = r.ch.PublishWithContext(ctx,
  "",     // exchange
  queue, // routing key
  false,  // mandatory
  false,  // immediate
  amqp091.Publishing {
    ContentType: "text/plain",
    Body:        []byte(body),
  })
  if err != nil {
	r.logger.Error().Err(err).Msgf("Failed to publish message")
	return err
  }
	return nil
}

func (r *rabbitMQService) Consume(ctx context.Context, queue string, handler func([]byte) error) error{
	_, err := r.ch.QueueDeclare(
		queue, // name
		true,    // durability
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		amqp091.Table{
		  amqp091.QueueTypeArg: amqp091.QueueTypeQuorum,
		},
	  )
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to declare queue %s", queue)
		return err
	}
	msgs, err := r.ch.Consume(
		queue, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to consume messages from queue %s", queue)
		return err
	}

	go func() {
		for  {
			select {
				case msgs, ok := <-msgs:
					if !ok {
						return
					}

					if err := handler(msgs.Body); err != nil {
						msgs.Nack(false, false)
					} else {
						msgs.Ack(false)
					}
				case <-ctx.Done():
					return
			}
		}
	}()
	return nil
}

func (r *rabbitMQService) Close() error {
	if r.ch != nil {
		if err := r.ch.Close(); err != nil {
			return err
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return err
		}
	}
 	return nil
}