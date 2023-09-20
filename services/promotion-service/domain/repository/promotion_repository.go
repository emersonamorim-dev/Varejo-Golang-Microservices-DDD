package repository

import (
	"Varejo-Golang-Microservices/services/promotion-service/domain/model"
	"Varejo-Golang-Microservices/services/promotion-service/infra/db"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPromotionRepository struct {
	client *mongo.Client
	kafka  *kafka.Producer
}

func NewMongoPromotionRepository(mongoURI string, kafkaBroker string) *MongoPromotionRepository {
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

	return &MongoPromotionRepository{
		client: client,
		kafka:  producer,
	}
}

func (r *MongoPromotionRepository) ListAll() ([]*model.Promotion, error) {
	promotionCollection := r.client.Database("promotionDB").Collection("promotions")

	cursor, err := promotionCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var promotions []*model.Promotion
	for cursor.Next(context.TODO()) {
		var promotion model.Promotion
		err := cursor.Decode(&promotion)
		if err != nil {
			return nil, err
		}
		promotions = append(promotions, &promotion)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return promotions, nil
}

func (r *MongoPromotionRepository) FindByID(id string) (*model.Promotion, error) {
	collection := r.client.Database("promotionDB").Collection("promotions")

	// Converte a string ID para ObjectID do MongoDB
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID da promoção inválido")
	}

	// Define o filtro para busca
	filter := bson.M{"_id": objID}

	var promotion model.Promotion
	err = collection.FindOne(context.TODO(), filter).Decode(&promotion)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("promoção não encontrada")
		}
		return nil, err
	}

	return &promotion, nil
}

func (r *MongoPromotionRepository) Save(promotion *model.Promotion) error {
	promotionCollection := r.client.Database("promotionDB").Collection("promotions")

	// Inserir a promoção na coleção
	_, err := promotionCollection.InsertOne(context.TODO(), promotion)
	if err != nil {
		log.Printf("Erro ao inserir promoção no MongoDB: %v\n", err)
		return err
	}

	// Converte a promoção para JSON para enviar ao Kafka
	promotionJSON, err := json.Marshal(promotion)
	if err != nil {
		log.Printf("Erro ao organizar a promoção: %v", err)
		return err
	}

	// Define o tópico do Kafka para a promoção
	topic := "Promotion_Topic_One"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          promotionJSON,
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

func (r *MongoPromotionRepository) Update(promotion *model.Promotion) error {
	promotionCollection := r.client.Database("promotionDB").Collection("promotions")

	// Usando o ID diretamente para o filtro
	filter := bson.M{"_id": promotion.ID}

	// Dados de atualização com todos os campos de Promotion
	updateData := bson.M{
		"title":         promotion.Title,
		"description":   promotion.Description,
		"startDate":     promotion.StartDate,
		"endDate":       promotion.EndDate,
		"discount":      promotion.Discount,
		"discountValue": promotion.DiscountValue,
		"status":        promotion.Status,
	}

	// Atualiza o documento
	_, err := promotionCollection.UpdateOne(context.TODO(), filter, bson.M{"$set": updateData})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoPromotionRepository) Delete(id string) error {
	// Seguindo o padrão do exemplo fornecido
	promotionCollection := r.client.Database("promotionDB").Collection("promotions")

	// Converte a string ID para um ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Define o filtro para encontrar a promoção pelo ID
	filter := bson.M{"_id": objID}

	// Deleta a promoção que corresponde ao filtro
	result, err := promotionCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	// Verifica se alguma promoção foi realmente deletada
	if result.DeletedCount == 0 {
		return errors.New("nenhuma promoção encontrada com o ID fornecido")
	}

	return nil
}

func (r *MongoPromotionRepository) Close() {
	r.kafka.Close()
}
