package main

import (
	"Varejo-Golang-Microservices/middleware"
	"Varejo-Golang-Microservices/routes"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI    = "mongodb://localhost:27017"
	kafkaBroker = "localhost:9092"
	serverPort  = ":8094"
)

func main() {
	r := setupRouter()
	r.Run(serverPort)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.TrustedProxyMiddleware())

	// Adicionado uma rota simples na raiz ("/")
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Bem-vindo à API Varejo Digital!",
		})
	})

	if err := setupMongoDB(); err != nil {
		log.Fatalf("Falha ao conectar ao MongoDB: %v", err)
	}

	//api_gateway.SetupGatewayRoutes(r)

	routes.SetupRoutes(r, mongoURI, kafkaBroker)
	return r
}

func setupMongoDB() error {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	return client.Ping(context.TODO(), nil)
}
