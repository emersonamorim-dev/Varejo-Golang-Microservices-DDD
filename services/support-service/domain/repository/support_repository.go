package repository

import (
	"Varejo-Golang-Microservices/services/support-service/domain/model"
	"Varejo-Golang-Microservices/services/support-service/infra/db"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSupportRepository struct {
	client *mongo.Client
	kafka  *kafka.Producer
}

func NewMongoSupportRepository(mongoURI string, kafkaBroker string) *MongoSupportRepository {
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

	return &MongoSupportRepository{
		client: client,
		kafka:  producer,
	}
}

// Lista Suporte
func (r *MongoSupportRepository) ListAll() ([]*model.Support, error) {
	supportCollection := r.client.Database("supportDB").Collection("supports")

	// Encontra todos os documentos na coleção
	cursor, err := supportCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Decodifica os documentos retornados para a estrutura Support
	var supports []*model.Support
	for cursor.Next(context.TODO()) {
		var support model.Support
		err := cursor.Decode(&support)
		if err != nil {
			continue
		}
		supports = append(supports, &support)
	}

	// Verifica se houve algum erro durante a iteração do cursor
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return supports, nil
}

func (r *MongoSupportRepository) FindByID(id string) (*model.Support, error) {
	collection := r.client.Database("supportDB").Collection("supports")

	// Converte a string ID para ObjectID do MongoDB
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID do suporte inválido")
	}

	// Define o filtro para busca
	filter := bson.M{"_id": objID}

	var support model.Support
	err = collection.FindOne(context.TODO(), filter).Decode(&support)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("suporte não encontrado")
		}
		return nil, err
	}

	return &support, nil
}

func (r *MongoSupportRepository) Save(support *model.Support) error {
	supportCollection := r.client.Database("supportDB").Collection("supports")

	// Insere o suporte na coleção
	_, err := supportCollection.InsertOne(context.TODO(), support)
	if err != nil {
		log.Printf("Erro ao inserir suporte no MongoDB: %v\n", err)
		return err
	}

	// Converte o suporte para JSON para enviar ao Kafka
	supportJSON, err := json.Marshal(support)
	if err != nil {
		log.Printf("Erro ao organizar o suporte: %v", err)
		return err
	}

	// Define o tópico do Kafka para o suporte
	topic := "Support_Topic_One"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          supportJSON,
	}

	// Envia a mensagem ao Kafka
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

func (r *MongoSupportRepository) Update(support *model.Support) error {
	supportCollection := r.client.Database("supportDB").Collection("supports")

	// Usando o ID diretamente para o filtro
	filter := bson.M{"_id": support.ID}

	// Dados de atualização com todos os campos de Support
	updateData := bson.M{
		"subject":     support.Subject,
		"message":     support.Message,
		"createdDate": support.CreatedDate,
		"data":        support.Data,
		"response":    support.Response,
		"status":      support.Status,
	}

	// Atualiza o documento
	_, err := supportCollection.UpdateOne(context.TODO(), filter, bson.M{"$set": updateData})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoSupportRepository) Delete(id string) error {
	supportCollection := r.client.Database("supportDB").Collection("supports")

	// Converte a string ID para um ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Define o filtro para encontrar o suporte pelo ID
	filter := bson.M{"_id": objID}

	// Deleta o suporte que corresponde ao filtro
	result, err := supportCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	// Verifica se algum suporte foi realmente deletado
	if result.DeletedCount == 0 {
		return errors.New("nenhum suporte encontrado com o ID fornecido")
	}

	return nil
}

// Fecha se tiver uma instância de Kafka
func (r *MongoSupportRepository) Close() {
	r.kafka.Close()
}
