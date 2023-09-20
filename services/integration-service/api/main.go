package main

import (
	"Varejo-Golang-Microservices/middleware"
	"Varejo-Golang-Microservices/services/integration-service/api/handler"
	"Varejo-Golang-Microservices/services/integration-service/domain/repository"
	"Varejo-Golang-Microservices/services/integration-service/domain/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	targetServiceURL   = "http://localhost:8099"
	defaultMongoURI    = "mongodb://localhost:27017"
	defaultKafkaBroker = "localhost:9092"
)

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

	// Crie uma instância do MongoIntegrationRepository
	mongoRepo := repository.NewMongoIntegrationRepository(mongoURI, kafkaBroker)

	// Use essa instância ao criar o IntegrationService
	integrationService := service.NewIntegrationService(mongoRepo)

	// Inicialize o manipulador (handler) com o serviço
	integrationHandler := handler.NewIntegrationHandler(integrationService)

	// Configurando as rotas
	authorized.GET("/integrate", integrationHandler.ListIntegrationData)
	authorized.GET("/integrate/:id", integrationHandler.GetIntegrationDataByID)
	authorized.POST("/integrate", integrationHandler.AddIntegrationData)
	authorized.PUT("/integrate/:id", integrationHandler.UpdateIntegrationData)
	authorized.DELETE("/integrate/:id", integrationHandler.DeleteIntegrationData)

	// Iniciando o servidor
	r.Run(":8082")
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
