package handler

import (
	"Varejo-Golang-Microservices/services/support-service/domain/model"
	"Varejo-Golang-Microservices/services/support-service/domain/service"
	"Varejo-Golang-Microservices/services/support-service/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SupportHandler struct {
	Service service.SupportService
}

// Inicializa um novo manipulador de suporte com o serviço fornecido
func NewSupportHandler(s service.SupportService) *SupportHandler {
	return &SupportHandler{
		Service: s,
	}
}

func convertDTOToSupport(supportDTO dto.SupportDTO) *model.Support {
	return &model.Support{
		ID:          primitive.NewObjectID(),
		Subject:     supportDTO.Subject,
		Message:     supportDTO.Message,
		CreatedDate: supportDTO.CreatedDate,
		Data:        supportDTO.Data,
		Response:    supportDTO.Response,
		Status:      supportDTO.Status,
	}
}

func convertDTOToSupportWithoutID(supportDTO dto.SupportDTO) *model.Support {
	return &model.Support{
		Subject:     supportDTO.Subject,
		Message:     supportDTO.Message,
		CreatedDate: supportDTO.CreatedDate,
		Response:    supportDTO.Response,
		Status:      supportDTO.Status,
		Category: model.Category{
			Name:        supportDTO.Category.Name,
			Description: supportDTO.Category.Description,
		},
	}
}

// Listar Suporte
func (h *SupportHandler) ListSupports(c *gin.Context) {
	supports, err := h.Service.ListAllSupports()
	if err != nil {
		log.Printf("Erro detalhado ao buscar suportes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar suportes"})
		return
	}

	if len(supports) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Nenhum suporte encontrado", "data": supports})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": supports})
}

// Busca Suporte por ID
func (h *SupportHandler) GetSupportByID(c *gin.Context) {
	supportID := c.Param("id")

	// Verifica se o ID do suporte foi fornecido
	if supportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID do suporte é obrigatório"})
		return
	}

	// Busca o suporte pelo ID
	support, err := h.Service.GetSupportByID(supportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar suporte"})
		return
	}

	// Verifica se o suporte foi encontrado
	if support == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Suporte não encontrado"})
		return
	}

	// Retorna o suporte encontrado
	c.JSON(http.StatusOK, support)
}

// Adiciona Suporte ao Banco de Dados
func (h *SupportHandler) AddSupport(c *gin.Context) {
	// Analisa os dados da solicitação na estrutura SupportDTO
	var supportDTO dto.SupportDTO
	if err := c.ShouldBindJSON(&supportDTO); err != nil {
		log.Printf("Erro ao ligar suporte JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar os dados do suporte."})
		return
	}

	// Converte DTO para a estrutura real de suporte
	support := convertDTOToSupport(supportDTO)

	// Salva o suporte usando o serviço
	err := h.Service.SaveSupport(support)
	if err != nil {
		log.Printf("Detalhes do Erro ao salvar suporte: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar suporte. Detalhes: " + err.Error()})
		return
	}

	// Atualizar o DTO com o ID gerado
	supportDTO.ID = support.ID.Hex()

	// Envia a resposta ao cliente
	c.JSON(http.StatusCreated, gin.H{"message": "Suporte cadastrado com sucesso.", "data": supportDTO})
}

func (h *SupportHandler) UpdateSupport(c *gin.Context) {
	// Obtem o ID do suporte da URL
	supportIDStr := c.Param("id")
	if supportIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Support ID is required"})
		return
	}

	// Converte a string ID para um ObjectID
	objID, err := primitive.ObjectIDFromHex(supportIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid support ID format"})
		return
	}

	// Vincula o JSON da solicitação ao DTO
	var supportDTO dto.SupportDTO
	if err := c.ShouldBindJSON(&supportDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converte o DTO para o modelo
	support := convertDTOToSupportWithoutID(supportDTO)

	// Define o ID do suporte
	support.ID = objID

	// Atualizar o suporte
	err = h.Service.UpdateSupport(support)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar o suporte. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Suporte atualizado com sucesso",
		"data":    supportDTO,
	})
}

func (h *SupportHandler) DeleteSupport(c *gin.Context) {
	supportID := c.Param("id")

	// Verifica se o ID do suporte foi fornecido
	if supportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID do suporte é obrigatório"})
		return
	}

	// Tenta deletar o suporte usando o serviço
	err := h.Service.DeleteSupport(supportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Erro ao excluir suporte. Detalhes: " + err.Error()})
		return
	}

	// Retorna uma mensagem de sucesso
	c.JSON(http.StatusOK, gin.H{"message": "Suporte deletado com sucesso"})
}
