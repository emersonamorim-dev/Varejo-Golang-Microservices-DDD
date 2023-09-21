package main

import (
	"Varejo-Golang-Microservices/middleware"
	"Varejo-Golang-Microservices/services/payment-service/api/handler"
	"Varejo-Golang-Microservices/services/payment-service/domain/repository"
	"Varejo-Golang-Microservices/services/payment-service/domain/service"
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

	// Inicialize as conexões de banco de dados, repositórios, serviços.
	paymentRepo := repository.NewMongoPaymentRepository(mongoURI, kafkaBroker)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Configurando as rotas
	r.GET("/payment", paymentHandler.GetAllPayments)
	r.GET("/payment/:id", paymentHandler.GetPaymentByID)
	r.POST("/payment", paymentHandler.AddPayment)
	r.PUT("/payment/:id", paymentHandler.UpdatePayment)
	r.DELETE("/payment/:id", paymentHandler.DeletePayment)

	// Starting the server
	r.Run(":8085")
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
