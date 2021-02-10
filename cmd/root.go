package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

var bootstrapServer string

// rootCmd represents the base command when called without any subcommands.
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

func kafkaClient() sarama.Client {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.V1_0_0_0
	kafkaConfig.Net.ReadTimeout = 5 * time.Second
	kafkaConfig.Net.DialTimeout = 5 * time.Second
	kafkaConfig.Net.WriteTimeout = 5 * time.Second
	addresses := []string{bootstrapServer}

	var client sarama.Client
	var err error

	for i := 0; i < 30; i++ {
		client, err = sarama.NewClient(addresses, kafkaConfig)
		if err == nil {
			break
		}

		fmt.Fprintln(os.Stderr, "failed to connect to "+bootstrapServer+" Retrying in 1s")
		time.Sleep(time.Second)
	}

	if err != nil {
		panic(err)
	}

	return client
}

func kafkaAdmin() sarama.ClusterAdmin {
	client := kafkaClient()

	var admin sarama.ClusterAdmin
	var err error

	for i := 0; i < 10; i++ {
		admin, err = sarama.NewClusterAdminFromClient(client)
		if err == nil {
			break
		}

		fmt.Fprintln(os.Stderr, "failed to admin cluster at "+bootstrapServer+" Retrying in 1s")
		time.Sleep(time.Second)
	}

	if err != nil {
		panic(err)
	}

	return admin
}
