package repository

import (
	"Varejo-Golang-Microservices/services/product-service/domain/model"
	"Varejo-Golang-Microservices/services/product-service/infra/db"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoProductRepository struct {
	client *mongo.Client
	kafka  *kafka.Producer
}

func NewMongoProductRepository(mongoURI string, kafkaBroker string) *MongoProductRepository {
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

	return &MongoProductRepository{
		client: client,
		kafka:  producer,
	}
}

func (r *MongoProductRepository) ListAll() ([]*model.Product, error) {
	productCollection := r.client.Database("productDB").Collection("products")
	cursor, err := productCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var products []*model.Product
	for cursor.Next(context.TODO()) {
		var product model.Product
		err := cursor.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *MongoProductRepository) FindByID(id string) (*model.Product, error) {
	collection := r.client.Database("productDB").Collection("products")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID do produto inválido")
	}

	filter := bson.M{"_id": objID}

	var product model.Product
	err = collection.FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("produto não encontrado")
		}
		return nil, err
	}

	return &product, nil
}

func (r *MongoProductRepository) SaveProduct(product *model.Product) error {
	productCollection := r.client.Database("productDB").Collection("products")

	// Inserir o produto na coleção
	_, err := productCollection.InsertOne(context.TODO(), product)
	if err != nil {
		log.Printf("Erro ao inserir produto no MongoDB: %v\n", err)
		return err
	}

	// Convertendo o produto para JSON para enviar ao Kafka
	productJSON, err := json.Marshal(product)
	if err != nil {
		log.Printf("Erro ao organizar o produto: %v", err)
		return err
	}

	// Definindo o tópico do Kafka para o produto
	topic := "Product_Topic_One"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          productJSON,
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

func (r *MongoProductRepository) Update(product *model.Product) error {
	productCollection := r.client.Database("productDB").Collection("products")

	// Usando o ID diretamente para o filtro
	filter := bson.M{"_id": product.ID}

	// Dados de atualização com todos os campos do Product
	updateData := bson.M{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"category": bson.M{
			"name":        product.Category.Name,
			"description": product.Category.Description,
		},
		"stock":     product.Stock,
		"addedDate": product.AddedDate,
	}

	// Atualiza o documento
	_, err := productCollection.UpdateOne(context.TODO(), filter, bson.M{"$set": updateData})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoProductRepository) Delete(id string) error {
	productCollection := r.client.Database("productDB").Collection("products")

	// Converte a string ID para um ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Define o filtro para encontrar o produto pelo ID
	filter := bson.M{"_id": objID}

	// Deleta o produto que corresponde ao filtro
	result, err := productCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	// Verifica se algum produto foi realmente deletado
	if result.DeletedCount == 0 {
		return errors.New("nenhum produto encontrado com o ID fornecido")
	}

	return nil
}

// Fecha a conexão Kafka
func (r *MongoProductRepository) Close() {
	r.kafka.Close()
}
