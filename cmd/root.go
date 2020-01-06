package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/Shopify/sarama"
)

var bootstrapServer string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kafkaCLI",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
