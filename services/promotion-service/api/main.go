package main

import (
	"Varejo-Golang-Microservices/middleware"
	"Varejo-Golang-Microservices/services/promotion-service/api/handler"
	"Varejo-Golang-Microservices/services/promotion-service/domain/repository"
	"Varejo-Golang-Microservices/services/promotion-service/domain/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const defaultMongoURI = "mongodb://localhost:27017"
const defaultKafkaBroker = "localhost:9092"

func main() {
	r := gin.Default()

	// Defina valores padrão para variáveis de ambiente.
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
	promotionRepo := repository.NewMongoPromotionRepository(mongoURI, kafkaBroker)
	promotionService := service.NewPromotionService(promotionRepo)

	// Initialize the handler with the services
	promotionHandler := handler.NewPromotionHandler(promotionService)

	// Setting up the routes
	r.GET("/promotions", promotionHandler.ListPromotions)
	r.GET("/promotions/:id", promotionHandler.GetPromotionByID)
	r.POST("/promotions", promotionHandler.AddPromotion)
	r.PUT("/promotions/:id", promotionHandler.UpdatePromotion)
	r.DELETE("/promotions/:id", promotionHandler.DeletePromotion)

	// Starting the server
	r.Run(":8087")
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
