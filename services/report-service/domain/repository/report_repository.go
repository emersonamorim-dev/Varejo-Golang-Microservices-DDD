package repository

import (
	"Varejo-Golang-Microservices/services/report-service/domain/model"
	"Varejo-Golang-Microservices/services/report-service/infra/db"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoReportRepository struct {
	client *mongo.Client
	kafka  *kafka.Producer
}

func NewMongoReportRepository(mongoURI string, kafkaBroker string) *MongoReportRepository {
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

	return &MongoReportRepository{
		client: client,
		kafka:  producer,
	}
}

func (r *MongoReportRepository) ListAll() ([]*model.Report, error) {
	reportCollection := r.client.Database("reportDB").Collection("reports")

	// Encontra todos os documentos na coleção
	cursor, err := reportCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Decodifica os documentos retornados para a estrutura Report
	var reports []*model.Report
	for cursor.Next(context.TODO()) {
		var report model.Report
		err := cursor.Decode(&report)
		if err != nil {
			continue
		}
		reports = append(reports, &report)
	}

	// Verifica se houve algum erro durante a iteração do cursor
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *MongoReportRepository) FindByID(id string) (*model.Report, error) {
	collection := r.client.Database("reportDB").Collection("reports")

	// Converte a string ID para ObjectID do MongoDB
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID do relatório inválido")
	}

	// Define o filtro para busca
	filter := bson.M{"_id": objID}

	var report model.Report
	err = collection.FindOne(context.TODO(), filter).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("relatório não encontrado")
		}
		return nil, err
	}

	return &report, nil
}

func (r *MongoReportRepository) Save(report *model.Report) error {
	reportCollection := r.client.Database("reportDB").Collection("reports")

	// Inserir o relatório na coleção
	_, err := reportCollection.InsertOne(context.TODO(), report)
	if err != nil {
		log.Printf("Erro ao inserir relatório no MongoDB: %v\n", err)
		return err
	}

	// Converte o relatório para JSON para enviar ao Kafka
	reportJSON, err := json.Marshal(report)
	if err != nil {
		log.Printf("Erro ao organizar o relatório: %v", err)
		return err
	}

	// Define o tópico do Kafka para o relatório
	topic := "Report_Topic_One"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          reportJSON,
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

func (r *MongoReportRepository) Update(report *model.Report) error {
	reportCollection := r.client.Database("reportDB").Collection("reports")

	// Usando o ID diretamente para o filtro
	filter := bson.M{"_id": report.ID}

	// Dados de atualização com todos os campos de Report
	updateData := bson.M{
		"title":       report.Title,
		"description": report.Description,
		"createdDate": report.CreatedDate,
		"data":        report.Data,
		"status":      report.Status,
	}

	// Atualiza o documento
	_, err := reportCollection.UpdateOne(context.TODO(), filter, bson.M{"$set": updateData})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoReportRepository) Delete(id string) error {
	reportCollection := r.client.Database("reportDB").Collection("reports")

	// Converte a string ID para um ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Define o filtro para encontrar o relatório pelo ID
	filter := bson.M{"_id": objID}

	// Deleta o relatório que corresponde ao filtro
	result, err := reportCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	// Verifica se algum relatório foi realmente deletado
	if result.DeletedCount == 0 {
		return errors.New("nenhum relatório encontrado com o ID fornecido")
	}

	return nil
}

// Fecha uma instância de Kafka
func (r *MongoReportRepository) Close() {
	r.kafka.Close()
}
