package cmd

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

// deleteTopicCmd represents the deleteTopic command.
var deleteTopicCmd = &cobra.Command{
	Use:   "deleteTopic",
	Short: "Delete one or more topics",
	Long: `deleteTopic deletes one or more topics. Example:

kafkaCLI deleteTopic --bootstrap-server kafka:9092 topic1 topic2
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		admin := kafkaAdmin()

		for _, topicName := range args {
			fmt.Println("Deleting topic " + topicName)
			err := admin.DeleteTopic(topicName)

			if err != nil {
				switch err {
				case sarama.ErrUnknownTopicOrPartition:
					fmt.Println("Topic " + topicName + " did not exist")
					break
				default:
					panic(err)
				}
			}
		}

		_ = admin.Close()
	},
}

func init() {
	rootCmd.AddCommand(deleteTopicCmd)
}
