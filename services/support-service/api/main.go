package main

import (
	"Varejo-Golang-Microservices/middleware"
	"Varejo-Golang-Microservices/services/support-service/api/handler"
	"Varejo-Golang-Microservices/services/support-service/domain/repository"
	"Varejo-Golang-Microservices/services/support-service/domain/service"
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
	supportRepo := repository.NewMongoSupportRepository(mongoURI, kafkaBroker)
	supportService := service.NewSupportService(supportRepo)
	supportHandler := handler.NewSupportHandler(supportService)

	// Configurando as rotas
	authorized.GET("/supports", supportHandler.ListSupports)
	authorized.GET("/supports/:id", supportHandler.GetSupportByID)
	authorized.POST("/supports", supportHandler.AddSupport)
	authorized.PUT("/supports/:id", supportHandler.UpdateSupport)
	authorized.DELETE("/supports/:id", supportHandler.DeleteSupport)

	// Starting the server
	r.Run(":8089")
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
