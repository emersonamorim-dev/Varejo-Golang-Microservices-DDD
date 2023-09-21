package handler

import (
	"Varejo-Golang-Microservices/services/order-service/domain/model"
	"Varejo-Golang-Microservices/services/order-service/domain/service"
	"Varejo-Golang-Microservices/services/order-service/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	Service service.OrderService
}

// Inicializa um novo manipulador de pedidos com o serviço fornecido
func NewOrderHandler(s service.OrderService) *OrderHandler {
	return &OrderHandler{
		Service: s,
	}
}

func convertDTOToOrder(dto dto.OrderDTO) model.Order {
	return model.Order{
		ID:              primitive.NewObjectID(),
		CustomerID:      dto.CustomerID,
		Products:        convertDTOItemsToOrderProducts(dto.Products),
		TotalPrice:      dto.TotalPrice,
		ShippingAddress: convertDTOAddressToModelAddress(dto.ShippingAddress),
		Status:          dto.Status,
		OrderDate:       dto.CreatedAt,
		DeliveryDate:    dto.UpdatedAt,
	}
}

func convertDTOAddressToModelAddress(address dto.Address) model.Address {
	return model.Address{
		Street:     address.Street,
		City:       address.City,
		State:      address.State,
		PostalCode: address.PostalCode,
		Country:    address.Country,
	}
}

func convertDTOItemsToOrderProducts(items []dto.OrderItemDTO) []model.OrderProduct {
	var products []model.OrderProduct

	for _, item := range items {
		product := model.OrderProduct{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price, // Use UnitPrice aqui
		}
		products = append(products, product)
	}

	return products
}

func convertDTOToOrderWithoutID(dto dto.OrderDTO) model.Order {
	products := convertDTOItemsToOrderProducts(dto.Products)
	return model.Order{
		CustomerID:   dto.CustomerID,
		Products:     products,
		TotalPrice:   dto.TotalPrice,
		Status:       dto.Status,
		OrderDate:    dto.CreatedAt,
		DeliveryDate: dto.UpdatedAt,
	}
}

// Listar Pedidos
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.Service.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar pedidos"})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Nenhum pedido encontrado", "data": orders})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(400, gin.H{"error": "O ID do pedido é obrigatório"})
		return
	}

	order, err := h.Service.GetOrderByID(orderID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao buscar pedido"})
		return
	}

	if order == nil {
		c.JSON(404, gin.H{"error": "Pedido não encontrado"})
		return
	}

	c.JSON(200, order)
}

func (h *OrderHandler) AddOrder(c *gin.Context) {
	// Analisa os dados da solicitação na estrutura OrderDTO
	var orderDTO dto.OrderDTO
	if err := c.ShouldBindJSON(&orderDTO); err != nil {
		log.Printf("Erro ao ligar pedido JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar os dados do pedido."})
		return
	}

	// Converti DTO para a estrutura de pedido real
	order := convertDTOToOrder(orderDTO)

	// Salva o pedido usando o serviço
	err := h.Service.SaveOrder(&order)
	if err != nil {
		log.Printf("Detalhes do Erro ao salvar pedido: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar pedido. Detalhes: " + err.Error()})
		return
	}

	// Envia a resposta ao cliente
	c.JSON(http.StatusOK, gin.H{"message": "Pedido cadastrado com sucesso.", "data": order})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("id")
	if orderIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	var orderDTO dto.OrderDTO
	if err := c.ShouldBindJSON(&orderDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converte orderDTO para um modelo de domínio
	order := convertDTOToOrderWithoutID(orderDTO)

	// Define o ID do pedido
	order.ID = objID

	// Atualiza apenas o status do pedido
	err = h.Service.UpdateOrderStatus(orderIDStr, order.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar pedido. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pedido atualizado com sucesso",
		"data":    orderDTO,
	})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID do pedido é obrigatório"})
		return
	}

	err := h.Service.DeleteOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Erro ao excluir pedido. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pedido deletado com sucesso"})
}
