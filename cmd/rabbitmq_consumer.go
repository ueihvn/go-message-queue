package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ueihvn/go-message-queue/config"
	"github.com/ueihvn/go-message-queue/pkg/rmqconnector"
	"go.uber.org/zap"
	"log"
)

var rabbitMQConsumer = &cobra.Command{
	Use:   "rmq_consumer",
	Short: "Starting RabbitMQ consumer",
	Long:  "Starting a simple RabbitMQ Consumer",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()
		log.Printf("rabbitMQConsumer start with config: %+v\n", conf)
		zapL, err := zap.NewDevelopment()
		logger := zapL.Sugar()
		if err != nil {
			log.Fatalf("get logger error: %+v", err)
		}
		deps := rmqconnector.IRMQConnectorDependencies{
			Config:      conf,
			Logger:      logger,
			ConsumeType: args[0],
		}
		var forever chan struct{}
		rmqConnector := rmqconnector.NewRMQConnector(deps)
		if err := rmqConnector.Consume(nil); err != nil {
			logger.Fatalf("RabbitMQ consumer consume err: %+v", err)
		}
		logger.Infof("consumer wait for message")
		<-forever
	},
}
