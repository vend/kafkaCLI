package cmd

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

// partitionerCmd finds partition associated to a given key
var partitionerCmd = &cobra.Command{
	Use:   "partitioner",
	Short: "Finds partition associated to a given key",
	Long: `createTopic can create one or more topics. Example:

kafkaCLI partitioner --topic blackhole --partition-count 16 --message test --key-id 00005a30-9766-11e3-a0f5-b8ca3a64f8f4
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		keyIdString := sarama.StringEncoder(keyID)
	
		part := sarama.NewHashPartitioner(topic)
		msg := &sarama.ProducerMessage{Topic: topic, Key: keyIdString, Value: sarama.StringEncoder(message)}
	
		partition, err := part.Partition(msg, int32(partitions))
		
		if err != nil {
			panic(err)
		}
	
		fmt.Printf("partition: %v key-id: %v\n", partition, keyIdString)

	},
}

var topic string
var partitionCount int16
var message string
var keyID string

func init() {
	rootCmd.AddCommand(partitionerCmd)

	partitionerCmd.Flags().StringVarP(&topic, "topic", "p", "blackhole", "topic name")
	partitionerCmd.Flags().Int16VarP(&partitionCount, "partition-count", "c", 16, "number of partitions for the topic")
	partitionerCmd.Flags().StringVarP(&message, "message", "m", "test", "message text")
	partitionerCmd.Flags().StringVarP(&keyID, "key-id", "k", "", "key id used by the message, e.g 00005a30-9766-11e3-a0f5-b8ca3a64f8f4")
}
