package cmd

import(
	"time"
)

var (
	partitions int32
	replicationFactor int16
	ifNotExists bool
	configEntries []string

	topic string
	partitionCount int16
	message string
	keyId string
	offsetType string
	messagesToRead int32
	bufferSize int32
	partitionsToRead string // [all] to read from all or comma separeted list
	consumptionDeadline time.Duration
)