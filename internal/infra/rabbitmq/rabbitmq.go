package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/datpham/user-service-ms/config"
	"github.com/datpham/user-service-ms/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *logger.Logger
	config  *config.Config
}

// NewRabbitMQClient creates a new RabbitMQ client
func NewRabbitMQClient(logger *logger.Logger, cfg *config.Config) (*RabbitMQ, error) {
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	conn, err := amqp.Dial(rabbitUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	if err = ch.ExchangeDeclare(
		cfg.RabbitMQ.ExchangeName, // exchange name
		"topic",                   // exchange type (topic allows routing based on wildcards)
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare an exchange: %v", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		logger:  logger,
		config:  cfg,
	}, nil
}

func (r *RabbitMQ) Setup(queueName, routingKey string) error {
	queue, err := r.DeclareQueue(queueName)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %s", err.Error())
	}

	if err := r.BindQueue(queue.Name, routingKey); err != nil {
		return fmt.Errorf("failed to bind queue: %s", err.Error())
	}

	return nil
}

// DeclareQueue declares a new queue
func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,  // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// BindQueue binds a queue to the exchange with a routing key
func (r *RabbitMQ) BindQueue(queueName, routingKey string) error {
	return r.channel.QueueBind(
		queueName,                      // queue name
		routingKey,                     // routing key
		r.config.RabbitMQ.ExchangeName, // exchange
		false,                          // no-wait
		nil,                            // arguments
	)
}

// Publish publishes a message to a queue
func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, message []byte) error {
	err := r.channel.PublishWithContext(
		ctx,
		r.config.RabbitMQ.ExchangeName, // exchange
		routingKey,                     // routing key
		false,                          // mandatory
		false,                          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	return nil
}

// Consume consumes messages from a queue
func (r *RabbitMQ) Consume(queueName, consumerName string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queueName,    // queue
		consumerName, // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // arguments
	)
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %v", err)
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %v", err)
		}
	}

	return nil
}
