package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"strings"
)

type TransactionService struct {
	repo     *repositories.TransactionRepository
	reportTZ string
}

func NewTransactionService(repo *repositories.TransactionRepository, reportTZ string) *TransactionService {
	if strings.TrimSpace(reportTZ) == "" {
		reportTZ = "Asia/Jakarta"
	}
	return &TransactionService{
		repo:     repo,
		reportTZ: reportTZ,
	}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items, s.reportTZ)
}

func (s *TransactionService) GetTransactions(page, limit int) ([]models.TransactionListItem, error) {
	return s.repo.GetTransactions(page, limit, s.reportTZ)
}

func (s *TransactionService) CountTransactions() (int, error) {
	return s.repo.CountTransactions()
}

func (s *TransactionService) SumTransactions() (int, error) {
	return s.repo.SumTransactions()
}

func (s *TransactionService) GetTransactionsByDateRange(page, limit int, startDate, endDate string) ([]models.TransactionListItem, error) {
	return s.repo.GetTransactionsByDateRange(page, limit, startDate, endDate, s.reportTZ)
}

func (s *TransactionService) CountTransactionsByDateRange(startDate, endDate string) (int, error) {
	return s.repo.CountTransactionsByDateRange(startDate, endDate, s.reportTZ)
}

func (s *TransactionService) SumTransactionsByDateRange(startDate, endDate string) (int, error) {
	return s.repo.SumTransactionsByDateRange(startDate, endDate, s.reportTZ)
}

func (s *TransactionService) GetTransactionByID(id int) (*models.Transaction, error) {
	return s.repo.GetTransactionByID(id, s.reportTZ)
}

func (s *TransactionService) GetTodayReport() (*models.DailyReport, error) {
	return s.repo.GetTodayReport(s.reportTZ)
}

func (s *TransactionService) GetReportRange(startDate, endDate string) (*models.DailyReport, error) {
	return s.repo.GetReportRange(startDate, endDate, s.reportTZ)
}
