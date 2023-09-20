package repository

import (
	"Varejo-Golang-Microservices/services/location-service/domain/model"
	"Varejo-Golang-Microservices/services/location-service/infra/db"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoLocationRepository struct {
	client *mongo.Client
	kafka  *kafka.Producer
}

func NewMongoLocationRepository(mongoURI string, kafkaBroker string) *MongoLocationRepository {
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

	return &MongoLocationRepository{
		client: client,
		kafka:  producer,
	}
}

func (r *MongoLocationRepository) FindByID(id string) (*model.Location, error) {
	collection := r.client.Database("locationDB").Collection("locations")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	var location model.Location
	err = collection.FindOne(context.TODO(), filter).Decode(&location)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &location, nil
}

func (r *MongoLocationRepository) ListAll() ([]*model.Location, error) {
	collection := r.client.Database("locationDB").Collection("locations")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var locations []*model.Location
	for cursor.Next(context.TODO()) {
		var location model.Location
		err := cursor.Decode(&location)
		if err != nil {
			return nil, err
		}
		locations = append(locations, &location)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *MongoLocationRepository) Save(location *model.Location) error {
	collection := r.client.Database("locationDB").Collection("locations")

	_, err := collection.InsertOne(context.TODO(), location)
	if err != nil {
		log.Printf("Erro ao inserir localização no MongoDB: %v\n", err)
		return err
	}

	// Convert location to JSON for Kafka
	locationJSON, err := json.Marshal(location)
	if err != nil {
		log.Printf("Erro ao organizar a localização: %v", err)
		return err
	}

	topic := "Location_Topic_One"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          locationJSON,
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

func (r *MongoLocationRepository) Update(location *model.Location) error {
	locationCollection := r.client.Database("locationDB").Collection("locations")

	// Criando o filtro usando bson.D para manter a consistência com o exemplo dado
	filter := bson.D{{Key: "_id", Value: location.ID}}

	// Estruturando os dados de atualização de acordo com o padrão fornecido
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "latitude", Value: location.Latitude},
			{Key: "longitude", Value: location.Longitude},
			{Key: "description", Value: location.Description},
			{Key: "address", Value: location.Address},
			{Key: "data", Value: location.Data},
			{Key: "createdDate", Value: location.CreatedDate},
			{Key: "status", Value: location.Status},
		}},
	}

	// Atualiza o documento
	_, err := locationCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *MongoLocationRepository) Delete(id string) error {
	locationCollection := r.client.Database("locationDB").Collection("locations")

	// Defini o filtro para a consulta usando o ID como uma string
	filter := bson.M{"_id": id}

	// Executa a operação de deleção
	result, err := locationCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	// Verifica se algum documento foi deletado
	if result.DeletedCount == 0 {
		return errors.New("nenhuma localização encontrada com o ID fornecido")
	}

	return nil
}
