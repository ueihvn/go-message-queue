package rmqconnector

import (
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

const (
	SimpleConsumer      = "simple"
	TaskQueueConsumer   = "task_queue"
	PubSubQueueConsumer = "pub_sub"
)

func (r *rmqConnectorImp) defaultConsumerCallBackFunc(msg amqp.Delivery) error {
	time.Sleep(1 * time.Second)
	r.logger.Infof("consumer recevice a message: %s", msg.Body)
	return nil
}

func (r *rmqConnectorImp) Consume(fnc ConsumerCallBackFunc) error {
	if fnc == nil {
		fnc = r.defaultConsumerCallBackFunc
	}
	switch r.consumeType {
	case SimpleConsumer:
		return r.simpleConsume(fnc)
	case TaskQueueConsumer:
		return r.taskQueueConsume(fnc)
	case PubSubQueueConsumer:
		return r.pubSubConsume(fnc)
	default:
		return fmt.Errorf("rmqConnectorImp produceType: %v not implemented", r.produceType)
	}
}

func (r *rmqConnectorImp) simpleConsume(fnc ConsumerCallBackFunc) error {
	queue, err := r.rmqChannel.QueueDeclare(
		SimpleConsumer, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("consume decleare queue error: %w", err)
	}

	msgs, err := r.rmqChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("consume register error: %w", err)
	}
	go func() {
		for msg := range msgs {
			if fnc == nil {
				r.logger.Infof("Consumer receive message: %+v", msg)
			} else {
				fnc(msg)
			}
		}
	}()
	return nil
}

func (r *rmqConnectorImp) taskQueueConsume(fnc ConsumerCallBackFunc) error {
	queue, err := r.rmqChannel.QueueDeclare(
		TaskQueueConsumer, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return fmt.Errorf("consume decleare queue error: %w", err)
	}

	if err := r.rmqChannel.Qos(
		1,
		0,
		false,
	); err != nil {
		return errors.Wrapf(err, "set QoS")
	}
	msgs, err := r.rmqChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("consume register error: %w", err)
	}
	go func() {
		for msg := range msgs {
			if fnc == nil {
				r.logger.Infof("Consumer receive message: %+v", msg)
			} else {
				fnc(msg)
				delay := time.Duration(len(msg.Body)) * time.Second
				time.Sleep(delay)
				r.logger.Infof("Done processing message in: %v", delay)
				msg.Ack(false)
			}
		}
	}()
	return nil
}

func (r *rmqConnectorImp) pubSubConsume(fnc ConsumerCallBackFunc) error {
	queue, err := r.rmqChannel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return errors.Wrap(err, "declare queue")
	}
	if err := r.rmqChannel.QueueBind(
		queue.Name,   // queue name
		"",           // routing key
		LogsExchange, // exchange
		false,
		nil,
	); err != nil {
		return errors.Wrap(err, "queue bind")
	}
	msgs, err := r.rmqChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return errors.Wrap(err, "consume register")
	}
	go func() {
		for msg := range msgs {
			if fnc == nil {
				r.logger.Infof("Consumer receive message: %+v", msg)
			} else {
				fnc(msg)
			}
		}
	}()
	return nil
}
