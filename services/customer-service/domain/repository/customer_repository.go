package repository

import (
	"Varejo-Golang-Microservices/services/customer-service/domain/model"
	"Varejo-Golang-Microservices/services/customer-service/infra/db"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCustomerRepository struct {
	client *mongo.Client
	kafka  *kafka.Producer
}

func NewMongoCustomerRepository(mongoURI string, kafkaBroker string) *MongoCustomerRepository {
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

	return &MongoCustomerRepository{
		client: client,
		kafka:  producer,
	}
}

func (m *MongoCustomerRepository) GetAll() ([]model.Customer, error) {
	collection := m.client.Database("customerDB").Collection("customers")

	// Busca todos os documentos na coleção.
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var customers []model.Customer
	// Itera sobre o cursor e decodifica cada documento em uma estrutura `Customer`.
	for cursor.Next(context.TODO()) {
		var customer model.Customer
		err := cursor.Decode(&customer)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (m *MongoCustomerRepository) FindByID(id string) (*model.Customer, error) {
	collection := m.client.Database("customerDB").Collection("customers")
	filter := bson.M{"_id": id}

	var customer model.Customer
	err := collection.FindOne(context.TODO(), filter).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("cliente não encontrado")
		}
		return nil, err
	}

	return &customer, nil
}

func (m *MongoCustomerRepository) Save(customer *model.Customer) error {
	collection := m.client.Database("customerDB").Collection("customers")

	_, err := collection.InsertOne(context.TODO(), customer)
	if err != nil {
		log.Printf("Erro ao inserir cliente no MongoDB: %v\n", err)
		return err
	}

	// Convert customer to JSON for Kafka
	customerJSON, err := json.Marshal(customer)
	if err != nil {
		log.Printf("Erro ao organizar o cliente: %v", err)
		return err
	}

	topic := "Customer_Topic_One"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          customerJSON,
	}

	deliveryChan := make(chan kafka.Event)
	err = m.kafka.Produce(message, deliveryChan)
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

func (m *MongoCustomerRepository) Update(customer *model.Customer) error {
	collection := m.client.Database("customerDB").Collection("customers")

	// Usando o ID diretamente para o filtro
	filter := bson.M{"_id": customer.ID}

	// Dados de atualização com todos os campos mencionados
	updateData := bson.M{
		"name":    customer.Name,
		"email":   customer.Email,
		"cell":    customer.Cell,
		"phone":   customer.Phone,
		"address": customer.Address,
		"zipCode": customer.ZipCode,
		"city":    customer.City,
	}

	// Atualiza o documento
	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": updateData})
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoCustomerRepository) Delete(id string) error {
	collection := m.client.Database("customerDB").Collection("customers")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("nenhum cliente encontrado com o ID fornecido")
	}

	return nil
}
