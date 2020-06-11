package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listTopicsCmd represents the listTopics command
var listTopicsCmd = &cobra.Command{
	Use:   "listTopics",
	Short: "List all the topics",
	Long: `listTopics all the topics. Example:

kafkaCLI listTopics --bootstrap-server kafka:9092
`,
	Run: func(cmd *cobra.Command, args []string) {
		admin := kafkaAdmin()
		topics, err := admin.ListTopics()
		if err != nil {
			panic(err)
		}

		for name, _ := range topics {
			fmt.Println(name)
		}

		_ = admin.Close()
	},
}

func init() {
	rootCmd.AddCommand(listTopicsCmd)
}
