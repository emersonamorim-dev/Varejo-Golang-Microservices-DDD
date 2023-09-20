package repository

import (
	"Varejo-Golang-Microservices/services/order-service/domain/model"
	"Varejo-Golang-Microservices/services/order-service/infra/db"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoOrderRepository struct {
	client *mongo.Client
	kafka  *kafka.Producer
}

func NewMongoOrderRepository(mongoURI string, kafkaBroker string) *MongoOrderRepository {
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

	return &MongoOrderRepository{
		client: client,
		kafka:  producer,
	}
}

func (r *MongoOrderRepository) GetAll() ([]*model.Order, error) {
	orderCollection := r.client.Database("orderDB").Collection("orders")
	cursor, err := orderCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var orders []*model.Order
	for cursor.Next(context.TODO()) {
		var order model.Order
		err := cursor.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *MongoOrderRepository) FindByID(id string) (*model.Order, error) {
	collection := r.client.Database("orderDB").Collection("orders")
	filter := bson.M{"_id": id}

	var order model.Order
	err := collection.FindOne(context.TODO(), filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("pedido não encontrado")
		}
		return nil, err
	}

	return &order, nil
}

func (r *MongoOrderRepository) Save(order *model.Order) error {
	orderCollection := r.client.Database("orderDB").Collection("orders")

	// Inseri o pedido na coleção
	_, err := orderCollection.InsertOne(context.TODO(), order)
	if err != nil {
		log.Printf("Erro ao inserir pedido no MongoDB: %v\n", err)
		return err
	}

	// Converte o pedido para JSON para enviar ao Kafka
	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Printf("Erro ao organizar o pedido: %v", err)
		return err
	}

	// Definindo o tópico do Kafka para o pedido
	topic := "Order_Topic_One"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          orderJSON,
	}

	// Enviando a mensagem ao Kafka
	deliveryChan := make(chan kafka.Event)
	err = r.kafka.Produce(message, deliveryChan)
	if err != nil {
		log.Printf("Erro ao produzir mensagem para Kafka: %v\n", err)
		close(deliveryChan)
		return err
	}

	// Manipulando a resposta do envio ao Kafka
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

func (r *MongoOrderRepository) Update(order *model.Order) error {
	collection := r.client.Database("orderDB").Collection("orders")

	// Usando o ID diretamente para o filtro
	filter := bson.M{"_id": order.ID}

	// Dados de atualização com todos os campos mencionados
	updateData := bson.M{
		"customerId":      order.CustomerID,
		"products":        order.Products,
		"totalPrice":      order.TotalPrice,
		"shippingAddress": order.ShippingAddress,
		"status":          order.Status,
		"orderDate":       order.OrderDate,
		"deliveryDate":    order.DeliveryDate,
	}

	// Atualiza o documento
	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": updateData})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoOrderRepository) Delete(id string) error {
	collection := r.client.Database("orderDB").Collection("orders")

	// Converte a string ID para um ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Define o filtro para encontrar o pedido pelo ID
	filter := bson.M{"_id": objID}

	// Deleta a ordem que corresponde ao filtro
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	// Verifica se algum pedido foi realmente deletado
	if result.DeletedCount == 0 {
		return errors.New("nenhuma ordem encontrada com o ID fornecido")
	}

	return nil
}

func (r *MongoOrderRepository) Close() {
	r.kafka.Close()
}
