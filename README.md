# kafkaCLI

The CLI tools that ship with Kafka are super slow, so I thought I'd implement
a subset of it in GO.

Currently only topic creation and deletion is supported.

```
kafkaCLI createTopic --bootstrap-server kafka:9092 --partitions 4 --replication-factor 1 --config message.format.version=2.0 --if-not-exists topic1 topic2
```

```
kafkaCLI deleteTopic --bootstrap-server kafka:9092 topic1 topic2
```