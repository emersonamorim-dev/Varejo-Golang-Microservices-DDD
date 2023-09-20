package main

import (
	"Varejo-Golang-Microservices/middleware"
	"Varejo-Golang-Microservices/services/customer-service/api/handler"
	"Varejo-Golang-Microservices/services/customer-service/domain/repository"
	"Varejo-Golang-Microservices/services/customer-service/domain/service"
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

	// Inicializar conexões de banco de dados, repositories, services, etc.
	customerRepo := repository.NewMongoCustomerRepository(mongoURI, kafkaBroker)
	customerService := service.NewCustomerService(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerService)

	// Define a rota de autenticação
	r.POST("/login", authenticate)

	// Grupo para endpoints protegidos por autenticação
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	{
		r.GET("/customers", customerHandler.GetAllCustomers)
		r.GET("/customers/:id", customerHandler.GetCustomerByID)
		r.POST("/customers", customerHandler.AddCustomer)
		r.PUT("/customers/:id", customerHandler.UpdateCustomer)
		r.DELETE("/customers/:id", customerHandler.DeleteCustomer)
	}

	r.Run(":8081")
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
