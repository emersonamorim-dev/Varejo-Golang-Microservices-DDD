package main

import (
	"Varejo-Golang-Microservices/middleware"
	"Varejo-Golang-Microservices/services/order-service/api/handler"
	"Varejo-Golang-Microservices/services/order-service/domain/repository"
	"Varejo-Golang-Microservices/services/order-service/domain/service"
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

	// Initialize the database connections, repositories e services.
	orderRepo := repository.NewMongoOrderRepository(mongoURI, kafkaBroker)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	// Setting up the routes
	r.GET("/order", orderHandler.GetAllOrders)
	r.GET("/orders/:id", orderHandler.GetOrderByID)
	r.POST("/order", orderHandler.AddOrder)
	r.PUT("/order/:id", orderHandler.UpdateOrderStatus)
	r.DELETE("/order/:id", orderHandler.DeleteOrder)

	// Starting the server
	r.Run(":8084")
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
