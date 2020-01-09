package cmd

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

// createTopicCmd represents the createTopic command
var createTopicCmd = &cobra.Command{
	Use:   "createTopic",
	Short: "Create one or more topics",
	Long: `createTopic can create one or more topics. Example:

kafkaCLI createTopic --bootstrap-server kafka:9092 --partitions 4 --replication-factor 1 --config message.format.version=2.0 --if-not-exists topic1 topic2
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := kafkaClient()
		if err != nil {
			panic(err)
		}
		kafkaAdmin, err := sarama.NewClusterAdminFromClient(client)
		if err != nil {
			panic(err)
		}

		for _, topicName := range args {
			fmt.Println("Creating topic " + topicName)
			err = kafkaAdmin.CreateTopic(topicName, topicDetail(), false)
			if err != nil {
				switch err.(type) {
				case *sarama.TopicError:
					if err, ok := err.(*sarama.TopicError); ok {
						if err.Err != sarama.ErrTopicAlreadyExists {
							panic(err)
						}

						if ifNotExists {
							fmt.Println(err)
						} else {
							panic(err)
						}
					}
					break
				default:
					panic(err)
				}
			}
		}

		_ = kafkaAdmin.Close()
	},
}

func init() {
	rootCmd.AddCommand(createTopicCmd)

	createTopicCmd.Flags().Int32VarP(&partitions, "partitions", "p", 1, "number of partitions for the topic")
	createTopicCmd.Flags().Int16VarP(&replicationFactor, "replication-factor", "r", 1, "replication-factor for the topic")
	createTopicCmd.Flags().BoolVarP(&ifNotExists, "if-not-exists", "", false, "only create the topic if it does not exist")
	createTopicCmd.Flags().StringSliceVarP(&configEntries, "config", "c", []string{}, "config")
}

func topicDetail() *sarama.TopicDetail {
	config := make(map[string]*string, 0)

	for _, entry := range configEntries {
		parts := strings.Split(entry, "=")
		config[parts[0]] = &parts[1]
	}

	detail := &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
		ConfigEntries:     config,
	}

	return detail
}
