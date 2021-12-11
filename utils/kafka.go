package utils

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"strconv"
	"sync"
)

const KafkaBroker = "localhost:9092"
const TopicName = "important"
const ClientID = "some_client"
const GroupConsumer = "some_group"
const KafkaVersion = "3.0.0"

func CheckKafka() {
	producer := createProducer([]string{KafkaBroker}, ClientID)
	for i:=0; i<10; i++{
		sendMessage(producer, TopicName, "Hello there"+strconv.Itoa(i))
	}

	group := createConsumerGroup(KafkaVersion, []string{KafkaBroker}, ClientID, GroupConsumer)
	ctx, cancel := context.WithCancel(context.Background())
	consumeMessages(ctx, group, TopicName)
	cancel()
	group.Close()

	closeProducer(producer)
}

func createProducer(brokerList []string, clientId string) sarama.SyncProducer {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	config.ClientID = clientId

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}

func sendMessage(producer sarama.SyncProducer, topic, message string) error {
	_, _, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	})

	return err
}

func closeProducer(producer sarama.SyncProducer) error {
	return producer.Close()
}

func createConsumerGroup(kafkaVersion string, brokerList []string, clientId string, groupName string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version, _ = sarama.ParseKafkaVersion(kafkaVersion)
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.ClientID = clientId

	consumerGroup, err := sarama.NewConsumerGroup(brokerList, groupName, config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v\n", err)
	}
	return consumerGroup
}

func consumeMessages(ctx context.Context, group sarama.ConsumerGroup, topic string) {
	consumer := Consumer{
		ready: make(chan bool),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := group.Consume(ctx, []string{topic}, &consumer); err != nil {
				log.Fatalf("Error from consumer: %v\n", err)
			}
			if ctx.Err() != nil {
				return
			}

			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	wg.Wait()
}

type Consumer struct {
	ready chan bool
}

func (consumer *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value=%s, timestamp=%v, topic=%s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}
