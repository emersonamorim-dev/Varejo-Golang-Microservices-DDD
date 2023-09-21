package main

import (
	"Varejo-Golang-Microservices/middleware"
	"Varejo-Golang-Microservices/services/location-service/api/handler"
	"Varejo-Golang-Microservices/services/location-service/domain/repository"
	"Varejo-Golang-Microservices/services/location-service/domain/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const defaultMongoURI = "mongodb://localhost:27017"
const defaultKafkaBroker = "localhost:9092"

func main() {
	r := gin.Default()

	// Define valores para variáveis de ambiente.
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = defaultMongoURI
	}

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = defaultKafkaBroker
	}

	r.POST("/login", authenticate)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	// Initialize database connections, repositories, services.
	locationRepo := repository.NewMongoLocationRepository(mongoURI, kafkaBroker)
	locationService := service.NewLocationService(locationRepo)
	locationHandler := handler.NewLocationHandler(locationService)

	// Configura as rotas
	r.GET("/locations", locationHandler.GetLocation)
	r.GET("/locations/:id", locationHandler.GetLocationByID)
	r.POST("/locations", locationHandler.AddLocation)
	r.PUT("/locations/:id", locationHandler.UpdateLocation)
	r.DELETE("/locations/:id", locationHandler.DeleteLocation)

	// Starting the server
	r.Run(":8083")
}

// rota de login
func authenticate(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "admin" && password == "password" {
		token, err := middleware.GenerateToken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao gerar token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token":   token,
			"message": "Credenciais válidas",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
}
