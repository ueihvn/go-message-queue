package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ueihvn/go-message-queue/config"
)

func Execute() error {
	rootCMD := &cobra.Command{
		Use:   "app",
		Short: "root command",
		Long:  "root command for starting application",
		Run:   func(_ *cobra.Command, args []string) {},
	}
	config.LoadConfig()

	rootCMD.AddCommand(rabbitMQProducer)
	rootCMD.AddCommand(rabbitMQConsumer)

	if err := rootCMD.Execute(); err != nil {
		panic(fmt.Errorf("execute root cmd, error: %w", err))
	}
	return nil
}
