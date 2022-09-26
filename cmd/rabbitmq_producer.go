package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ueihvn/go-message-queue/config"
	"log"
)

var rabbitMQProducer = &cobra.Command{
	Use:   "rmq_producer",
	Short: "Starting RabbitMQ producer",
	Long:  "Statinng a simple RabbitMQ Producer",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()
		log.Printf("rabbitMQProducer start with config: %+v\n", conf)
	},
}
