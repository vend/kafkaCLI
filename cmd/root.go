package cmd

import (
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

var bootstrapServer string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kafkaCLI",
	Short: "A CLI tool for Kafka",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&bootstrapServer, "bootstrap-server", "s", "", "address of a node in the kafka cluster")
	_ = rootCmd.MarkPersistentFlagRequired("bootstrap-server")
}

func kafkaClient() (sarama.Client, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.V1_0_0_0
	addresses := []string{bootstrapServer}
	return sarama.NewClient(addresses, kafkaConfig)
}
