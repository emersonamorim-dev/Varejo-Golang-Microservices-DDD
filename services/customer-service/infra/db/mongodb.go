package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DefaultMongoURI = "mongodb://localhost:27017"

func ConnectMongoDB(connectionString string) (*mongo.Client, error) {
	if connectionString == "" {
		connectionString = DefaultMongoURI
	}

	clientOptions := options.Client().ApplyURI(connectionString)

	// Define um timeout de 10 segundos para o contexto
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Pinga o banco de dados para garantir que a conexão foi estabelecida
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Falha ao conectar-se a MongoDB: %v", err)
		return nil, err
	}

	log.Println("Conectado com sucesso a MongoDB")
	return client, nil
}
