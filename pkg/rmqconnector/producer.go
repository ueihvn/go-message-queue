package rmqconnector

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

const (
	SimpleProducer      = "simple"
	TaskQueueProducer   = "task_queue"
	PubSubQueueProducer = "pub_sub"
	LogsExchange        = "logs"
)

type rmqConnectorImp struct {
	rmqConnection *amqp.Connection
	rmqChannel    *amqp.Channel
	logger        *zap.SugaredLogger
	produceType   string
	consumeType   string
}

func NewRMQConnector(deps IRMQConnectorDependencies) IRMQConnector {
	rmqURI := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		deps.Config.RabbitMQ.UserName,
		deps.Config.RabbitMQ.Password,
		deps.Config.RabbitMQ.Host,
		deps.Config.RabbitMQ.Port,
	)
	conn, err := amqp.Dial(rmqURI)
	if err != nil {
		deps.Logger.Fatalf("connect to RabbitMQ error: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		deps.Logger.Fatalf("open chanel to RabbitMQ error: %v", err)
	}
	return &rmqConnectorImp{
		rmqConnection: conn,
		rmqChannel:    ch,
		logger:        deps.Logger,
		produceType:   deps.ProduceType,
		consumeType:   deps.ConsumeType,
	}
}

func (r *rmqConnectorImp) Produce(messageBody []byte) error {
	switch r.produceType {
	case SimpleProducer:
		return r.simpleProduce(messageBody)
	case TaskQueueProducer:
		return r.taskQueueProduce(messageBody)
	case PubSubQueueProducer:
		return r.pubSubProduce(messageBody)
	default:
		return fmt.Errorf("rmqConnectorImp produceType: %v not implemented", r.produceType)
	}
}

func (r *rmqConnectorImp) simpleProduce(messageBody []byte) error {
	queue, err := r.rmqChannel.QueueDeclare(
		SimpleProducer, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("declare a queue error: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.rmqChannel.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        messageBody,
		},
	); err != nil {
		return fmt.Errorf("rmq produce message error: %w", err)
	}
	return nil
}

func (r *rmqConnectorImp) taskQueueProduce(messageBody []byte) error {
	queue, err := r.rmqChannel.QueueDeclare(
		TaskQueueProducer, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return errors.Wrapf(err, "declare a queue")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.rmqChannel.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         messageBody,
		},
	); err != nil {
		return errors.Wrapf(err, "rmq produce message")
	}
	return nil
}

func (r *rmqConnectorImp) pubSubProduce(messageBody []byte) error {
	if err := r.rmqChannel.ExchangeDeclare(
		LogsExchange,        // name
		amqp.ExchangeFanout, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	); err != nil {
		return errors.Wrap(err, "declare exchange")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.rmqChannel.PublishWithContext(
		ctx,
		LogsExchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        messageBody,
		},
	); err != nil {
		return fmt.Errorf("rmq produce message error: %w", err)
	}
	return nil
}
