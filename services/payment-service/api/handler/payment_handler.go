package handler

import (
	"Varejo-Golang-Microservices/services/payment-service/domain/model"
	"Varejo-Golang-Microservices/services/payment-service/domain/service"
	"Varejo-Golang-Microservices/services/payment-service/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentHandler struct {
	Service service.PaymentService
}

// AuthMiddleware para autenticação e autorização
func (h *PaymentHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// Inicializa um novo manipulador de pagamento com o serviço fornecido
func NewPaymentHandler(s service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		Service: s,
	}
}

// Função auxiliar para converter PaymentDTO em model.Payment
func convertDTOPayment(paymentDTO dto.PaymentDTO) (model.Payment, error) {
	objID, err := primitive.ObjectIDFromHex(paymentDTO.ID)
	if err != nil {
		return model.Payment{}, err
	}

	return model.Payment{
		ID:          objID, // Aqui é feita a conversão da string para ObjectID
		OrderID:     paymentDTO.OrderID,
		CustomerID:  paymentDTO.CustomerID,
		Amount:      paymentDTO.Amount,
		Method:      convertDTOPaymentMethod(paymentDTO.Method),
		Status:      paymentDTO.Status,
		PaymentDate: paymentDTO.PaymentDate,
	}, nil
}

// Função auxiliar para converter PaymentMethodDTO em model.PaymentMethod
func convertDTOPaymentMethod(methodDTO dto.PaymentMethodDTO) model.PaymentMethod {
	return model.PaymentMethod{
		Type:       methodDTO.Type,
		CardNumber: methodDTO.CardNumber,
		Expiry:     methodDTO.Expiry,
		CVV:        methodDTO.CVV,
	}
}

// Função auxiliar para converter PaymentDTO em model.Payment, mas sem o ID.
func convertDTOPaymentWithoutID(paymentDTO dto.PaymentDTO) model.Payment {
	return model.Payment{
		Amount:      paymentDTO.Amount,
		Method:      convertDTOPaymentMethod(paymentDTO.Method),
		Status:      paymentDTO.Status,
		PaymentDate: paymentDTO.PaymentDate,
	}
}

func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
	payments, err := h.Service.GetAllPayments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar pagamentos"})
		return
	}

	if len(payments) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Nenhum pagamento encontrado", "data": payments})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(400, gin.H{"error": "O ID de pagamento é obrigatório"})
		return
	}

	payment, err := h.Service.GetPaymentByID(paymentID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao buscar pagamento"})
		return
	}

	if payment == nil {
		c.JSON(404, gin.H{"error": "Pagamento não encontrado"})
		return
	}

	c.JSON(200, payment)
}

func (h *PaymentHandler) AddPayment(c *gin.Context) {
	// Verifica os dados da solicitação na estrutura PaymentDTO
	var paymentDTO dto.PaymentDTO
	if err := c.ShouldBindJSON(&paymentDTO); err != nil {
		log.Printf("Erro ao ligar pagamento JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar os dados do pagamento."})
		return
	}

	// Converte DTO para a estrutura de pagamento real
	payment, err := convertDTOPayment(paymentDTO)
	if err != nil {
		log.Printf("Erro ao converter DTO: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar o ID do pagamento."})
		return
	}

	// Salva o pagamento usando o serviço
	err = h.Service.SavePayment(&payment)
	if err != nil {
		log.Printf("Detalhes do Erro ao salvar pagamento: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar pagamento. Detalhes: " + err.Error()})
		return
	}

	// Envia a resposta ao cliente
	c.JSON(http.StatusOK, gin.H{"message": "Pagamento efetuado com sucesso.", "data": payment})
}

// Atualizar Clientes por ID
func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
	paymentIDStr := c.Param("id")
	if paymentIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}

	// Converta a string do ID para o tipo ObjectID
	objID, err := primitive.ObjectIDFromHex(paymentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID format"})
		return
	}

	var paymentDTO dto.PaymentDTO
	if err := c.ShouldBindJSON(&paymentDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converta customerDTO para um modelo de domínio
	payment := convertDTOPaymentWithoutID(paymentDTO)

	// Defina o ID do cliente
	payment.ID = objID

	err = h.Service.UpdatePayment(&payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar Pagamento. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cliente atualizado com sucesso",
		"data":    paymentDTO,
	})
}

// Deletar Pagamento por ID
func (h *PaymentHandler) DeletePayment(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID do Pagamento é obrigatório"})
		return
	}

	err := h.Service.DeletePayment(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Erro ao excluir pagamento. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pagamento deletado com sucesso"})
}
