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

	consumerCmd.Flags().StringVarP(&topic, "topic", "t", "blackhole", "topic name")
	consumerCmd.Flags().StringVarP(&partitionsToRead, "partitions-to-read", "p", "all", "comma list partitions to read or all to read from all partitions")
	consumerCmd.Flags().StringVarP(&offsetType, "offset-type", "n", "newest", "read from latest or newest offset?")
	consumerCmd.Flags().Int32VarP(&messagesToRead, "messages-to-read", "m", 5000, "number of messages to read")
	consumerCmd.Flags().Int32VarP(&bufferSize, "buffer-size", "b", 256, "buffer size of the message channel.")
	consumerCmd.Flags().Int32VarP(&consumptionDeadline, "consumption-deadline", "d", 3*time.Minute, "stop reading messages after X minutes")
	consumerCmd.MarkFlagRequired("key-id")

}
