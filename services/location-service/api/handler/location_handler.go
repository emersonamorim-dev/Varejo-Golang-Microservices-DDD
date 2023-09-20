package handler

import (
	"Varejo-Golang-Microservices/services/location-service/domain/model"
	"Varejo-Golang-Microservices/services/location-service/domain/service"
	"Varejo-Golang-Microservices/services/location-service/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LocationHandler struct {
	Service service.LocationService
}

// Inicializa um novo manipulador de localização com o serviço fornecido
func NewLocationHandler(s service.LocationService) *LocationHandler {
	return &LocationHandler{
		Service: s,
	}
}

// Função para ajudar a converter o método Update:
func convertDTOToLocationDataWithoutID(locationDTO dto.LocationDTO) model.Location {
	return model.Location{
		Latitude:    locationDTO.Latitude,
		Longitude:   locationDTO.Longitude,
		Description: locationDTO.Description,
		Address:     locationDTO.Address,
		Data:        locationDTO.Data,
		CreatedDate: locationDTO.CreatedDate,
		Status:      locationDTO.Status,
	}
}

func (h *LocationHandler) GetLocation(c *gin.Context) {
	locationID := c.Query("id")

	// Se o ID não for fornecido, retorne todas as localizações
	if locationID == "" {
		locations, err := h.Service.GetAllLocations()
		if err != nil {
			c.JSON(500, gin.H{"error": "Erro ao buscar todas as localizações"})
			return
		}

		if len(locations) == 0 {
			c.JSON(404, gin.H{"error": "Nenhuma localização encontrada"})
			return
		}

		c.JSON(200, locations)
		return
	}

	location, err := h.Service.GetLocationByID(locationID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao buscar localização"})
		return
	}

	if location == nil {
		c.JSON(404, gin.H{"error": "Localização não encontrada"})
		return
	}

	c.JSON(200, location)
}

func (h *LocationHandler) GetLocationByID(c *gin.Context) {
	locationID := c.Param("id")
	if locationID == "" {
		c.JSON(400, gin.H{"error": "O ID da localização é obrigatório"})
		return
	}

	location, err := h.Service.GetLocationByID(locationID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao buscar localização"})
		return
	}

	if location == nil {
		c.JSON(404, gin.H{"error": "Localização não encontrada"})
		return
	}

	c.JSON(200, location)
}

func (h *LocationHandler) AddLocation(c *gin.Context) {
	var locationDTO dto.LocationDTO
	if err := c.ShouldBindJSON(&locationDTO); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Usei UUID para gerar um ID único para a localização.
	uniqueID := uuid.New().String()

	// Converte locationDTO para model.Location
	location := &model.Location{
		ID:          uniqueID,
		Latitude:    locationDTO.Latitude,
		Longitude:   locationDTO.Longitude,
		Description: locationDTO.Description,
		Address:     locationDTO.Address,
		Data:        locationDTO.Data,
		CreatedDate: locationDTO.CreatedDate,
		Status:      locationDTO.Status,
	}

	err := h.Service.SaveLocation(location)
	if err != nil {
		c.JSON(500, gin.H{"error": "Erro ao adicionar localização"})
		return
	}

	// Construir a resposta com a localização e a mensagem
	response := gin.H{
		"location": location,
		"message":  "Localização e Tópico adicionada com sucesso!",
	}

	c.JSON(201, response)
}

func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	// Obtive o ID da localização do parâmetro da rota
	locationIDStr := c.Param("id")
	if locationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID da localização é obrigatório"})
		return
	}

	// Liguei o corpo JSON da solicitação ao LocationDTO
	var locationDTO dto.LocationDTO
	if err := c.ShouldBindJSON(&locationDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converti locationDTO para um modelo de domínio
	locationData := convertDTOToLocationDataWithoutID(locationDTO)

	// Defini o ID do LocationData
	locationData.ID = locationIDStr

	// Atualiza o LocationData usando o serviço
	err := h.Service.UpdateLocation(&locationData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar a localização. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Localização atualizada com sucesso",
		"data":    locationDTO,
	})
}

func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	locationID := c.Param("id")
	if locationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID da localização é obrigatório"})
		return
	}

	// Tenta deletar a localização usando o serviço
	err := h.Service.DeleteLocation(locationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Erro ao excluir localização. Detalhes: " + err.Error()})
		return
	}

	// Retorna sucesso se a localização foi deletada
	c.JSON(http.StatusOK, gin.H{"message": "Localização deletada com sucesso"})
}
