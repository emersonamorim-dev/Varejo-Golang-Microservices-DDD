package handler

import (
	"Varejo-Golang-Microservices/services/customer-service/domain/model"
	"Varejo-Golang-Microservices/services/customer-service/domain/service"
	"Varejo-Golang-Microservices/services/customer-service/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerHandler struct {
	Service service.CustomerService
}

// AuthMiddleware para autenticação e autorização
func (h *CustomerHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func NewCustomerHandler(s service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		Service: s,
	}
}

func convertDTOToCustomer(dto dto.CustomerDTO) model.Customer {
	return model.Customer{
		ID:      primitive.NewObjectID(),
		Name:    dto.Name,
		Email:   dto.Email,
		Cell:    dto.Cell,
		Phone:   dto.Phone,
		Address: dto.Address,
		ZipCode: dto.ZipCode,
		City:    dto.City,
	}
}

// Função de conversão sem considerar o ID
func convertDTOToCustomerWithoutID(dto dto.CustomerDTO) model.Customer {
	return model.Customer{
		Name:    dto.Name,
		Email:   dto.Email,
		Cell:    dto.Cell,
		Phone:   dto.Phone,
		Address: dto.Address,
		ZipCode: dto.ZipCode,
		City:    dto.City,
	}
}

// Listar Clientes
func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	customers, err := h.Service.GetAllCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar clientes"})
		return
	}

	if len(customers) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Nenhum cliente encontrado", "data": customers})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// Listar Clientes por id
func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	customerID := c.Param("id")
	if customerID == "" {
		c.JSON(400, gin.H{"error": "O ID do cliente é obrigatório"})
		return
	}

	customer, err := h.Service.GetCustomerByID(customerID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao buscar cliente"})
		return
	}

	if customer == nil {
		c.JSON(404, gin.H{"error": "Cliente não encontrado"})
		return
	}

	c.JSON(200, customer)
}

// Salva Clientes
func (h *CustomerHandler) AddCustomer(c *gin.Context) {
	var customerDTO dto.CustomerDTO
	if err := c.ShouldBindJSON(&customerDTO); err != nil {
		log.Printf("Erro ao ligar cliente JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar os dados do cliente."})
		return
	}

	customer := convertDTOToCustomer(customerDTO)

	err := h.Service.SaveCustomer(&customer)
	if err != nil {
		log.Printf("Detalhes do Erro ao salvar cliente: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar cliente. Detalhes: " + err.Error()})
		return
	}

	// Enviar a resposta para o cliente
	c.JSON(http.StatusOK, gin.H{"message": "Cliente cadastrado com sucesso e Tópico criado com sucesso."})
}

// Atualizar Clientes por ID
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	customerIDStr := c.Param("id")
	if customerIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	// Converta a string do ID para o tipo ObjectID
	objID, err := primitive.ObjectIDFromHex(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID format"})
		return
	}

	var customerDTO dto.CustomerDTO
	if err := c.ShouldBindJSON(&customerDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converta customerDTO para um modelo de domínio
	customer := convertDTOToCustomerWithoutID(customerDTO)

	// Defina o ID do cliente
	customer.ID = objID

	err = h.Service.UpdateCustomer(&customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar cliente. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cliente atualizado com sucesso",
		"data":    customerDTO,
	})
}

// Deletar Clientes por ID
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	customerID := c.Param("id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID do cliente é obrigatório"})
		return
	}

	err := h.Service.DeleteCustomer(customerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Erro ao excluir cliente. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cliente deletado com sucesso"})
}
