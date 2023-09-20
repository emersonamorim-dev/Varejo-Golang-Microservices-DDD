package service

import (
	"Varejo-Golang-Microservices/services/report-service/domain/model"
	"Varejo-Golang-Microservices/services/report-service/domain/repository"
)

type ReportService interface {
	ListAllReports() ([]*model.Report, error)
	GetReportByID(id string) (*model.Report, error)
	SaveReport(report *model.Report) error
	UpdateReport(report *model.Report) error
	DeleteReport(id string) error
}

type ReportServiceImpl struct {
	reportRepo *repository.MongoReportRepository
}

func NewReportService(reportRepo *repository.MongoReportRepository) ReportService {
	return &ReportServiceImpl{
		reportRepo: reportRepo,
	}
}

func (s *ReportServiceImpl) ListAllReports() ([]*model.Report, error) {
	return s.reportRepo.ListAll()
}

func (s *ReportServiceImpl) GetReportByID(id string) (*model.Report, error) {
	return s.reportRepo.FindByID(id)
}

func (s *ReportServiceImpl) SaveReport(report *model.Report) error {
	return s.reportRepo.Save(report)
}

func (s *ReportServiceImpl) UpdateReport(report *model.Report) error {
	return s.reportRepo.Update(report)
}

func (s *ReportServiceImpl) DeleteReport(id string) error {
	return s.reportRepo.Delete(id)
}
