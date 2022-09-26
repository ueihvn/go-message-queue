package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ueihvn/go-message-queue/config"
	"log"
)

var rabbitMQConsumer = &cobra.Command{
	Use:   "rmq_consumer",
	Short: "Starting RabbitMQ consumer",
	Long:  "Statinng a simple RabbitMQ Consumer",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()
		log.Printf("rabbitMQConsumer start with config: %+v\n", conf)
	},
}
