package cmd

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"time"

	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
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

		kafkaClient := kafkaClient()

		consumerGroup, err := sarama.NewConsumer([]string{bootstrapServer}, kafkaClient.Config())
		if err != nil {
			panic(err)
		}
	
		var initialOffset int64
		switch offsetType {
		case "oldest":
			initialOffset = sarama.OffsetOldest
		case "newest":
			initialOffset = sarama.OffsetNewest
		default:
			//Unkown offsetType. Defaulting to Newest"
			initialOffset = sarama.OffsetNewest
		}
	
		var (
			messages   = make(chan *sarama.ConsumerMessage, bufferSize)
			closing    = make(chan struct{})
			ErrTimeout = errors.New("Timed out waiting")
			wg       sync.WaitGroup
		)
	
		go func() {
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Kill, os.Interrupt)
			<-signals
			log.Info("Initiating shutdown of consumer...")
			close(closing)
		}()
	
		partitionList, err := getPartitions(consumerGroup)
		if err != nil {
			log.WithError(err).Fatal("Failed to get the list of partitions")

		}
	
		for _, partition := range partitionList {
			pc, err := consumerGroup.ConsumePartition(topic, partition, initialOffset)
			log.WithFields(log.Fields{
				"partition":  partition,
				"topic": topic,
			}).Info("consumer started")			
			if err != nil {
				log.WithFields(log.Fields{
					"partition":  partition,
				}).WithError(err).Fatal("Failer to consume messages")
			}
	
			go func(pc sarama.PartitionConsumer) {
				<-closing
				pc.AsyncClose()
			}(pc)
	
			wg.Add(1)
			go func(pc sarama.PartitionConsumer) {
				defer wg.Done()
				// for message := range pc.Messages() {
				// 	messages <- message
				// }
				for {
					select {
					case message := <-pc.Messages():
						if message == nil {
							// We have no more messages apparently
							log.Warn("We have no more messages apparently... exiting")
							return
						}
						messages <- message
					case err := <-pc.Errors():
						log.WithError(err).Error("Error consuming message")
					case <-time.After(consumptionDeadline):
						log.WithFields(log.Fields{
							"timeout":  consumptionDeadline,
						}).WithError(ErrTimeout).Error("consumption deadline reached")
						return
					}
				}
	
			}(pc)
		}
	
		go func() {
			messageCountStart := 0
			for msg := range messages {
				messageCountStart++
				log.WithFields(log.Fields{
					"partition":  msg.Partition,
					"offset": msg.Offset,
					"key": string(msg.Key),
					"value": string(msg.Value),
				}).Info("message received")
				if messageCountStart > messagesToRead {
					log.WithFields(log.Fields{
						"count":  messageCountStart,
					}).Info("Reached messages threshold")
					close(closing)
				}
			}
		}()
	
		wg.Wait()
		log.WithFields(log.Fields{
			"topic":  topic,
		}).Info("Done consuming topic")
	
		close(messages)
	
		if consumerGroup != nil {
			if err := consumerGroup.Close(); err != nil {
				log.WithError(err).Error("Failed to close consumer group")
			}
		}

	},
}


func getPartitions(c sarama.Consumer) ([]int32, error) {
	if partitionsToRead == "all" {
		return c.Partitions(topic)
	}

	tmp := strings.Split(partitionsToRead, ",")
	var pList []int32
	for i := range tmp {
		val, err := strconv.ParseInt(tmp[i], 10, 32)
		if err != nil {
			return nil, err
		}
		pList = append(pList, int32(val))
	}

	return pList, nil
}

func init() {
	rootCmd.AddCommand(consumerCmd)

	consumerCmd.Flags().StringVarP(&topic, "topic", "t", "blackhole", "topic name")
	consumerCmd.Flags().StringVarP(&partitionsToRead, "partitions-to-read", "p", "all", "comma list partitions to read or all to read from all partitions")
	consumerCmd.Flags().StringVarP(&offsetType, "offset-type", "n", "newest", "read from latest or newest offset?")
	consumerCmd.Flags().IntVarP(&messagesToRead, "messages-to-read", "m", 5000, "number of messages to read")
	consumerCmd.Flags().Int32VarP(&bufferSize, "buffer-size", "b", 256, "buffer size of the message channel.")
	consumerCmd.Flags().DurationVarP(&consumptionDeadline, "consumption-deadline", "d", 3*time.Minute, "stop reading messages after X minutes")
	consumerCmd.MarkFlagRequired("key-id")

}
