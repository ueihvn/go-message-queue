package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ueihvn/go-message-queue/config"
	"github.com/ueihvn/go-message-queue/pkg/rmqconnector"
	"go.uber.org/zap"
	"log"
	"time"
)

var rabbitMQProducer = &cobra.Command{
	Use:   "rmq_producer",
	Short: "Starting RabbitMQ producer",
	Long:  "Statinng a simple RabbitMQ Producer",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()
		log.Printf("rabbitMQProducer start with config: %+v\n", conf)
		zapL, err := zap.NewDevelopment(zap.AddStacktrace(zap.ErrorLevel))
		logger := zapL.Sugar()
		if err != nil {
			log.Fatalf("get logger error: %+v", err)
		}
		deps := rmqconnector.IRMQConnectorDependencies{
			Config:      conf,
			Logger:      logger,
			ProduceType: args[0],
		}
		rmqConnector := rmqconnector.NewRMQConnector(deps)
		indexMessage := 0
		for {
			time.Sleep(1 * time.Second)
			messageBody := []byte(fmt.Sprintf("produce message index %d", indexMessage))
			if err := rmqConnector.Produce(messageBody); err != nil {
				logger.Errorf("rmqConnector error produce message: %s, error: %+v", messageBody, err)
			} else {
				logger.Infof("rmqConnector successfully produce message index: %d", indexMessage)
			}
			indexMessage++
		}
	},
}
