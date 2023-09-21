package handler

import (
	"Varejo-Golang-Microservices/services/integration-service/domain/model"
	"Varejo-Golang-Microservices/services/integration-service/domain/service"
	"Varejo-Golang-Microservices/services/integration-service/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IntegrationHandler struct {
	Service service.IntegrationService
}

// AuthMiddleware para autenticação e autorização
func (h *IntegrationHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// Inicializa um novo manipulador de integração com o serviço fornecido
func NewIntegrationHandler(s service.IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{
		Service: s,
	}
}

// Função de conversão considerando o ID
func convertDTOToIntegrationData(dto dto.IntegrationDTO) model.IntegrationData {
	objectID, err := primitive.ObjectIDFromHex(dto.ID)
	if err != nil {
		// Aqui você deve decidir o que fazer em caso de erro.
		objectID = primitive.NewObjectID()
	}
	return model.IntegrationData{
		ID:       objectID,
		Name:     dto.Name,
		Endpoint: dto.Endpoint,
		APIKey:   dto.APIKey,
		Data:     dto.Data,
		Other:    dto.Other,
	}
}

// Função de conversão sem considerar o ID
func convertDTOToIntegrationDataWithoutID(dto dto.IntegrationDTO) model.IntegrationData {
	return model.IntegrationData{
		Name:     dto.Name,
		Endpoint: dto.Endpoint,
		APIKey:   dto.APIKey,
		Data:     dto.Data,
		Other:    dto.Other,
	}
}

func (h *IntegrationHandler) GetIntegrationDataByID(c *gin.Context) {
	dataID := c.Param("id")
	if dataID == "" {
		c.JSON(400, gin.H{"error": "O ID de integração é obrigatório"})
		return
	}

	data, err := h.Service.GetIntegrationDataByID(dataID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao buscar dados de integração"})
		return
	}

	if data == nil {
		c.JSON(404, gin.H{"error": "Dados de integração não encontrados"})
		return
	}

	c.JSON(200, data)
}

// Salva Integrações
func (h *IntegrationHandler) AddIntegrationData(c *gin.Context) {
	var dataDTO dto.IntegrationDTO
	if err := c.ShouldBindJSON(&dataDTO); err != nil {
		log.Printf("Erro ao ligar IntegrationDTO JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar os dados de integração."})
		return
	}

	// Converte dataDTO para model.IntegrationData sem considerar o ID
	data := &model.IntegrationData{
		Name:     dataDTO.Name,
		Endpoint: dataDTO.Endpoint,
		APIKey:   dataDTO.APIKey,
		Data:     dataDTO.Data,
		Other:    dataDTO.Other,
	}

	// Declarando a variável err aqui
	var err error
	err = h.Service.SaveIntegrationData(data)
	if err != nil {
		log.Printf("Detalhes do Erro ao salvar dados de integração: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar dados de integração. Detalhes: " + err.Error()})
		return
	}

	// Enviar a resposta para o cliente
	c.JSON(http.StatusOK, gin.H{"message": "Dados de integração adicionados com sucesso."})
}

func (h *IntegrationHandler) UpdateIntegrationData(c *gin.Context) {
	// Obtive o ID da integração do parâmetro da rota
	integrationIDStr := c.Param("id")
	if integrationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Integration ID is required"})
		return
	}

	// Converti a string do ID para o tipo ObjectID
	objID, err := primitive.ObjectIDFromHex(integrationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid integration ID format"})
		return
	}

	// Liguei o corpo JSON da solicitação ao IntegrationDTO
	var dataDTO dto.IntegrationDTO
	if err := c.ShouldBindJSON(&dataDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converti dataDTO para um modelo de domínio
	data := convertDTOToIntegrationDataWithoutID(dataDTO)

	// Defini o ID do IntegrationData
	data.ID = objID

	// Atualizei o IntegrationData usando o serviço
	err = h.Service.UpdateIntegrationData(&data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar os dados de integração. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dados de integração atualizados com sucesso",
		"data":    dataDTO,
	})
}

func (h *IntegrationHandler) DeleteIntegrationData(c *gin.Context) {
	dataID := c.Param("id")
	if dataID == "" {
		c.JSON(400, gin.H{"error": "O ID de integração é obrigatório"})
		return
	}

	err := h.Service.DeleteIntegrationData(dataID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao excluir os dados de integração"})
		return
	}

	c.Status(200)
}

// Buscar todos os dados de integração do banco de dados
func (h *IntegrationHandler) ListIntegrationData(c *gin.Context) {
	datas, err := h.Service.ListAllIntegrationData()
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao buscar dados de integração"})
		return
	}

	c.JSON(200, datas)
}
