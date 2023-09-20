package handler

import (
	"Varejo-Golang-Microservices/services/promotion-service/domain/model"
	"Varejo-Golang-Microservices/services/promotion-service/domain/service"
	"Varejo-Golang-Microservices/services/promotion-service/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromotionHandler struct {
	Service service.PromotionService
}

// Inicializa um novo manipulador de promoção com o serviço fornecido
func NewPromotionHandler(s service.PromotionService) *PromotionHandler {
	return &PromotionHandler{
		Service: s,
	}
}

func convertDTOToPromotion(promotionDTO dto.PromotionDTO) *model.Promotion {
	return &model.Promotion{
		ID:            primitive.NewObjectID(),
		Title:         promotionDTO.Title,
		Description:   promotionDTO.Description,
		StartDate:     promotionDTO.StartDate,
		EndDate:       promotionDTO.EndDate,
		Discount:      promotionDTO.Discount,
		DiscountValue: promotionDTO.DiscountValue,
		Status:        promotionDTO.Status,
	}
}

func convertDTOToPromotionWithoutID(promotionDTO dto.PromotionDTO) *model.Promotion {
	return &model.Promotion{
		Title:         promotionDTO.Title,
		Description:   promotionDTO.Description,
		StartDate:     promotionDTO.StartDate,
		EndDate:       promotionDTO.EndDate,
		Discount:      promotionDTO.Discount,
		DiscountValue: promotionDTO.DiscountValue,
		Status:        promotionDTO.Status,
	}
}

// Listar Promoções
func (h *PromotionHandler) ListPromotions(c *gin.Context) {
	promotions, err := h.Service.ListAllPromotions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar promoções"})
		return
	}

	if len(promotions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Nenhuma promoção encontrada", "data": promotions})
		return
	}

	c.JSON(http.StatusOK, promotions)
}

func (h *PromotionHandler) GetPromotionByID(c *gin.Context) {
	promotionID := c.Param("id")

	// Verifica se o ID da promoção foi fornecido
	if promotionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID da promoção é obrigatório"})
		return
	}

	// Busca a promoção pelo ID
	promotion, err := h.Service.GetPromotionByID(promotionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar promoção"})
		return
	}

	// Verifica se a promoção foi encontrada
	if promotion == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Promoção não encontrada"})
		return
	}

	// Retorna a promoção encontrada
	c.JSON(http.StatusOK, promotion)
}

func (h *PromotionHandler) AddPromotion(c *gin.Context) {
	// Analisa os dados da solicitação na estrutura PromotionDTO
	var promotionDTO dto.PromotionDTO
	if err := c.ShouldBindJSON(&promotionDTO); err != nil {
		log.Printf("Erro ao ligar promoção JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar os dados da promoção."})
		return
	}

	// Converte DTO para a estrutura real de promoção
	promotion := convertDTOToPromotion(promotionDTO)

	// Salva a promoção usando o serviço
	err := h.Service.SavePromotion(promotion)
	if err != nil {
		log.Printf("Detalhes do Erro ao salvar promoção: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar promoção. Detalhes: " + err.Error()})
		return
	}

	// Envia a resposta ao cliente
	c.JSON(http.StatusCreated, gin.H{"message": "Promoção cadastrada com sucesso.", "data": promotionDTO})
}

func (h *PromotionHandler) UpdatePromotion(c *gin.Context) {
	promotionIDStr := c.Param("id")
	if promotionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Promotion ID is required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(promotionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid promotion ID format"})
		return
	}

	var promotionDTO dto.PromotionDTO
	if err := c.ShouldBindJSON(&promotionDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converte promotionDTO para um modelo de domínio
	promotion := convertDTOToPromotionWithoutID(promotionDTO)

	// Define o ID da promoção
	promotion.ID = objID

	// Atualiza a promoção
	err = h.Service.UpdatePromotion(promotion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar a promoção. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Promoção atualizada com sucesso",
		"data":    promotionDTO,
	})
}

func (h *PromotionHandler) DeletePromotion(c *gin.Context) {
	promotionID := c.Param("id")

	// Verifica se o ID da promoção foi fornecido
	if promotionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID da promoção é obrigatório"})
		return
	}

	// Tenta deletar a promoção usando o serviço
	err := h.Service.DeletePromotion(promotionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Erro ao excluir promoção. Detalhes: " + err.Error()})
		return
	}

	// Retorna uma mensagem de sucesso
	c.JSON(http.StatusOK, gin.H{"message": "Promoção deletada com sucesso"})
}
