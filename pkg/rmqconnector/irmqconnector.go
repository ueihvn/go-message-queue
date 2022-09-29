package rmqconnector

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ueihvn/go-message-queue/config"
	"go.uber.org/zap"
)

type IRMQConnector interface {
	Produce(messageBody []byte) error
	Consume(backFunc ConsumerCallBackFunc) error
}

type IRMQConnectorDependencies struct {
	Config      *config.Config
	Logger      *zap.SugaredLogger
	ProduceType string
	ConsumeType string
}

type ConsumerCallBackFunc func(amqp.Delivery) error
