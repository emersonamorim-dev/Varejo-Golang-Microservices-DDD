package repository

import (
	"Varejo-Golang-Microservices/services/payment-service/domain/model"
	"Varejo-Golang-Microservices/services/payment-service/infra/db"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPaymentRepository struct {
	client *mongo.Client
	kafka  *kafka.Producer
}

func NewMongoPaymentRepository(mongoURI string, kafkaBroker string) *MongoPaymentRepository {
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

	return &MongoPaymentRepository{
		client: client,
		kafka:  producer,
	}
}

func (r *MongoPaymentRepository) GetAllPayments() ([]*model.Payment, error) {
	paymentCollection := r.client.Database("paymentDB").Collection("payments")

	// Buscando todos os pagamentos
	cursor, err := paymentCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var payments []*model.Payment
	for cursor.Next(context.TODO()) {
		var payment model.Payment
		err := cursor.Decode(&payment)
		if err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *MongoPaymentRepository) FindByID(id string) (*model.Payment, error) {
	collection := r.client.Database("paymentDB").Collection("payments")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID inválido")
	}
	filter := bson.M{"_id": objID}

	var payment model.Payment

	err = collection.FindOne(context.TODO(), filter).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("pagamento não encontrado")
		}
		return nil, err
	}

	return &payment, nil
}

func (r *MongoPaymentRepository) Save(payment *model.Payment) error {
	collection := r.client.Database("paymentDB").Collection("payments")

	// Insere o pagamento na coleção
	_, err := collection.InsertOne(context.TODO(), payment)
	if err != nil {
		log.Printf("Erro ao inserir pagamento no MongoDB: %v\n", err)
		return err
	}

	// Converte o pagamento para JSON para enviar ao Kafka
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		log.Printf("Erro ao organizar o pagamento: %v", err)
		return err
	}

	// Definindo o tópico do Kafka para o pagamento
	topic := "Payment_Topic_One"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          paymentJSON,
	}

	// Enviando a mensagem ao Kafka
	deliveryChan := make(chan kafka.Event)
	err = r.kafka.Produce(message, deliveryChan)
	if err != nil {
		log.Printf("Erro ao produzir mensagem para Kafka: %v\n", err)
		close(deliveryChan)
		return err
	}

	// Manipula a resposta do envio ao Kafka
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

func (r *MongoPaymentRepository) Update(payment *model.Payment) error {
	collection := r.client.Database("paymentDB").Collection("payments")

	// Usei o ID para filtrar
	filter := bson.M{"_id": payment.ID}

	// Dados de atualização, mencionando todos os campos no objeto Payment
	updateData := bson.M{
		"amount":      payment.Amount,
		"method":      payment.Method,
		"status":      payment.Status,
		"paymentDate": payment.PaymentDate,
	}

	// Atualiza o documento
	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": updateData})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoPaymentRepository) Delete(id string) error {
	collection := r.client.Database("paymentDB").Collection("payments")

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
		return errors.New("nenhuma ordem encontrada com o ID fornecido")
	}

	return nil
}
