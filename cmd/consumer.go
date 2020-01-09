package cmd

import (
	//"fmt"
	"time"
	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

// consumerCmd consumes message from kafka
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consumes [x] messages from kafka and exits after [consumption-deadkube]",
	Long: `Consumes messages from kafka. Example:

kafkaCLI consumer --topic blackhole --partitions-to-read all --offset-type newest --messages-to-read 5000 --buffer-size 256
--consumption-deadline 120
`,
	Run: func(cmd *cobra.Command, args []string) {



	},
}


func init() {
	rootCmd.AddCommand(consumerCmd)

	consumerCmd.Flags().StringVarP(&topic, "topic", "p", "blackhole", "topic name")
	consumerCmd.Flags().StringVarP(&partitionsToRead, "partitions-to-read", "c", 16, "number of partitions for the topic")
	consumerCmd.Flags().StringVarP(&message, "message", "m", "test", "message text")
	consumerCmd.Flags().StringVarP(&keyId, "key-id", "k", "", "key id used by the message, e.g 00005a30-9766-11e3-a0f5-b8ca3a64f8f4")
	consumerCmd.MarkFlagRequired("key-id")

}
