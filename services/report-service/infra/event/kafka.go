package event

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ProduceMessage envia uma mensagem para um tópico do Kafka.
func ProduceMessage(broker, topic, message string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return err
	}

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	p.Close()
	return nil
}

// ConsumeMessage consome mensagens de um tópico do Kafka.
func ConsumeMessage(broker, topic, groupID string, messageChan chan string) error {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return err
	}

	c.Subscribe(topic, nil)

	for {
		msg, err := c.ReadMessage(100 * time.Millisecond)
		if err == nil {
			messageChan <- string(msg.Value)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
