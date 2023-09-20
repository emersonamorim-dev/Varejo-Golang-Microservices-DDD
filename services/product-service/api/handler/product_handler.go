package handler

import (
	"Varejo-Golang-Microservices/services/product-service/domain/model"
	"Varejo-Golang-Microservices/services/product-service/domain/service"
	"Varejo-Golang-Microservices/services/product-service/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	Service service.ProductService
}

// Inicializa um novo manipulador de produto com o serviço fornecido
func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{
		Service: s,
	}
}

func convertDTOToProduct(dto dto.ProductDTO) *model.Product {
	return &model.Product{
		ID:          primitive.NewObjectID(),
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Category: model.Category{
			Name:        dto.Category.Name,
			Description: dto.Category.Description,
		},
		Stock:     dto.Stock,
		AddedDate: dto.AddedDate,
		Status:    dto.Status,
	}
}

func convertDTOToProductWithoutID(dto dto.ProductDTO) *model.Product {
	return &model.Product{
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Category: model.Category{
			Name:        dto.Category.Name,
			Description: dto.Category.Description,
		},
	}
}

// Listar Produtos
func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.Service.ListAllProducts()
	if err != nil {
		log.Printf("Erro ao buscar relatórios: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar relatórios"})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Nenhum produto encontrado", "data": products})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Busca um produto por seu ID.
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID := c.Param("id")

	// Verifica se o ID do produto foi fornecido
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID do produto é obrigatório"})
		return
	}

	// Busca o produto pelo ID
	product, err := h.Service.GetProductByID(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar produto"})
		return
	}

	// Verifica se o produto foi encontrado
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
		return
	}

	// Retorna o produto encontrado
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) AddProduct(c *gin.Context) {
	// Analisa os dados da solicitação na estrutura ProductDTO
	var productDTO dto.ProductDTO
	if err := c.ShouldBindJSON(&productDTO); err != nil {
		log.Printf("Erro ao ligar produto JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar os dados do produto."})
		return
	}

	// Converte DTO para a estrutura real de produto
	product := convertDTOToProduct(productDTO)

	// Salva o produto usando o serviço
	err := h.Service.SaveProduct(product)
	if err != nil {
		log.Printf("Detalhes do Erro ao salvar produto: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar produto. Detalhes: " + err.Error()})
		return
	}

	// Envia a resposta ao cliente
	c.JSON(http.StatusCreated, gin.H{"message": "Produto cadastrado com sucesso.", "data": productDTO})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	var productDTO dto.ProductDTO
	if err := c.ShouldBindJSON(&productDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converte productDTO para um modelo de domínio
	product := convertDTOToProductWithoutID(productDTO)

	// Define o ID do produto
	product.ID = objID

	// Atualiza o produto
	err = h.Service.UpdateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar o produto. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produto atualizado com sucesso",
		"data":    productDTO,
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID do produto é obrigatório"})
		return
	}

	err := h.Service.DeleteProduct(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Erro ao excluir produto. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produto deletado com sucesso"})
}
