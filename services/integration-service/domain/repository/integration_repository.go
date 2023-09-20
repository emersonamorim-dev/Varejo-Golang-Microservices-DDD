package repository

import (
	"Varejo-Golang-Microservices/services/integration-service/domain/model"
	"Varejo-Golang-Microservices/services/integration-service/infra/db"
	"context"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoIntegrationRepository struct {
	client *mongo.Client
	kafka  *kafka.Producer
}

func NewMongoIntegrationRepository(mongoURI string, kafkaBroker string) *MongoIntegrationRepository {
	client, err := db.ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("Erro ao conectar-se ao MongoDB: %v", err)
	}

	config := &kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Erro ao conectar-se ao Kafka: %v", err)
	}

	return &MongoIntegrationRepository{
		client: client,
		kafka:  producer,
	}
}

func (r *MongoIntegrationRepository) ListAllIntegrationData() ([]*model.IntegrationData, error) {
	integrationCollection := r.client.Database("integrationDB").Collection("integrations")

	cursor, err := integrationCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var integrations []*model.IntegrationData
	err = cursor.All(context.TODO(), &integrations)
	if err != nil {
		return nil, err
	}
	return integrations, nil
}

func (r *MongoIntegrationRepository) FindByID(id string) (*model.IntegrationData, error) {
	integrationCollection := r.client.Database("integrationDB").Collection("integrations")

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	var integration model.IntegrationData

	err := integrationCollection.FindOne(context.TODO(), filter).Decode(&integration)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &integration, nil
}

func (r *MongoIntegrationRepository) Save(data *model.IntegrationData) error {
	integrationCollection := r.client.Database("integrationDB").Collection("integrations")

	_, err := integrationCollection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Printf("Erro ao inserir dados de integração no MongoDB: %v\n", err)
		return err
	}

	// Convert integration data to JSON for Kafka
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Printf("Erro ao organizar os dados de integração: %v", err)
		return err
	}

	topic := "Integration_Topic_One"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          dataJSON,
	}

	deliveryChan := make(chan kafka.Event)
	err = r.kafka.Produce(message, deliveryChan)
	if err != nil {
		log.Printf("Erro ao produzir mensagem para Kafka: %v\n", err)
		close(deliveryChan)
		return err
	}

	e := <-deliveryChan
	switch ev := e.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			log.Printf("Erro ao enviar a mensagem ao Kafka: %v\n", ev.TopicPartition.Error)
			close(deliveryChan)
			return ev.TopicPartition.Error
		}
	}

	close(deliveryChan)
	return nil
}

func (r *MongoIntegrationRepository) Update(data *model.IntegrationData) error {
	integrationCollection := r.client.Database("integrationDB").Collection("integrations")

	filter := bson.D{{Key: "_id", Value: data.ID}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: data.Name},
			{Key: "endpoint", Value: data.Endpoint},
			{Key: "api_key", Value: data.APIKey},
			{Key: "data", Value: data.Data},
			{Key: "other", Value: data.Other},
		}},
	}

	_, err := integrationCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *MongoIntegrationRepository) DeleteIntegrationData(id string) error {
	integrationCollection := r.client.Database("integrationDB").Collection("integrations")

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	_, err := integrationCollection.DeleteOne(context.TODO(), filter)
	return err
}
