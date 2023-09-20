package handler

import (
	"Varejo-Golang-Microservices/services/report-service/domain/model"
	"Varejo-Golang-Microservices/services/report-service/domain/service"
	"Varejo-Golang-Microservices/services/report-service/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportHandler struct {
	Service service.ReportService
}

// Inicializa um novo manipulador de relatórios com o serviço fornecido
func NewReportHandler(s service.ReportService) *ReportHandler {
	return &ReportHandler{
		Service: s,
	}
}

func convertDTOToReport(reportDTO dto.ReportDTO) *model.Report {
	return &model.Report{
		ID:          primitive.NewObjectID(),
		Title:       reportDTO.Title,
		Description: reportDTO.Description,
		CreatedDate: reportDTO.CreatedDate,
		Data:        reportDTO.Data,
		Status:      model.ReportStatus(reportDTO.Status),
	}
}

func convertDTOToReportWithoutID(dto dto.ReportDTO) *model.Report {
	return &model.Report{
		Title:       dto.Title,
		Description: dto.Description,
		CreatedDate: dto.CreatedDate,
		Data:        dto.Data,
		Status:      dto.Status,
		Category: model.Category{
			Name:        dto.Category.Name,
			Description: dto.Category.Description,
		},
	}
}

// Listar Relatórios
func (h *ReportHandler) ListReports(c *gin.Context) {
	reports, err := h.Service.ListAllReports()
	if err != nil {
		log.Printf("Erro detalhado ao buscar relatórios: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar relatórios"})
		return
	}

	if len(reports) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Nenhum relatório encontrado", "data": reports})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reports})
}

func (h *ReportHandler) GetReportByID(c *gin.Context) {
	reportID := c.Param("id")

	// Verifica se o ID do relatório foi fornecido
	if reportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID do relatório é obrigatório"})
		return
	}

	// Busca o relatório pelo ID
	report, err := h.Service.GetReportByID(reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar relatório"})
		return
	}

	// Verifica se o relatório foi encontrado
	if report == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Relatório não encontrado"})
		return
	}

	// Retorna o relatório encontrado
	c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) AddReport(c *gin.Context) {
	// Analisa os dados da solicitação na estrutura ReportDTO
	var reportDTO dto.ReportDTO
	if err := c.ShouldBindJSON(&reportDTO); err != nil {
		log.Printf("Erro ao ligar relatório JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar os dados do relatório."})
		return
	}

	// Converte DTO para a estrutura real de relatório
	report := convertDTOToReport(reportDTO)

	// Salva o relatório usando o serviço
	err := h.Service.SaveReport(report)
	if err != nil {
		log.Printf("Detalhes do Erro ao salvar relatório: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar relatório. Detalhes: " + err.Error()})
		return
	}

	// Atualizar o DTO com o ID gerado
	reportDTO.ID = report.ID.Hex()

	// Envia a resposta ao cliente
	c.JSON(http.StatusCreated, gin.H{"message": "Relatório cadastrado com sucesso.", "data": reportDTO})
}

func (h *ReportHandler) UpdateReport(c *gin.Context) {
	reportIDStr := c.Param("id")
	if reportIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID is required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(reportIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID format"})
		return
	}

	var reportDTO dto.ReportDTO
	if err := c.ShouldBindJSON(&reportDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converter reportDTO em modelo de domínio
	report := convertDTOToReportWithoutID(reportDTO)

	// Defina o ID do relatório
	report.ID = objID

	// Atualizar o relatório
	err = h.Service.UpdateReport(report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar o relatório. Detalhes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Relatório atualizado com sucesso",
		"data":    reportDTO,
	})
}

func (h *ReportHandler) DeleteReport(c *gin.Context) {
	reportID := c.Param("id")

	// Verifica se o ID do relatório foi fornecido
	if reportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O ID do relatório é obrigatório"})
		return
	}

	// Tenta deletar o relatório usando o serviço
	err := h.Service.DeleteReport(reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Erro ao excluir relatório. Detalhes: " + err.Error()})
		return
	}

	// Retorna uma mensagem de sucesso
	c.JSON(http.StatusOK, gin.H{"message": "Relatório deletado com sucesso"})
}
