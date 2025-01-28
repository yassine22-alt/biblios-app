package service

import (
	"bookstore/api/api/internal/model"
	"bookstore/api/api/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

type ReportService struct {
	OrderRepo repository.OrderStore
	BookRepo  repository.BookStore
}

func NewReportService(orderRepo repository.OrderStore, bookRepo repository.BookStore) *ReportService {
	return &(ReportService{OrderRepo: orderRepo,
		BookRepo: bookRepo})
}

func (s *ReportService) StartSalesReportGenrator(ctx context.Context, logger *log.Logger) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Println("Generating sales report...")
			if err := s.GenerateReport(ctx); err != nil {
				logger.Printf("Error generating sales report: %v\n", err)
			} else {
				logger.Println("Sales report generated successfully.")

			}

		}
	}
}

func (s *ReportService) SaveReportAsJSON(report model.ReportModel) error {
	const reportDir = "./reports"
	if err := os.MkdirAll(reportDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create reports directory: %v", err)
	}

	timestamp := report.GeneratedAt.Format("20060102_150405")
	filePath := fmt.Sprintf("%s/report_%s.json", reportDir, timestamp)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("failed to write report to file: %v", err)
	}

	log.Printf("Report saved to %s\n", filePath)
	return nil
}

func (s *ReportService) GenerateReport(ctx context.Context) error {
	m := make(map[string]string)
	orders, err := s.OrderRepo.SearchOrders(ctx, m)

	if err != nil {
		return err
	}

	report := model.ReportModel{}
	report.TotalOrders = s.TotalOrders(ctx, orders)
	report.TotalRevenue = s.TotalRevenue(ctx, orders)
	report.TotalBooksSold = s.TotalBooksSold(ctx, orders)
	report.TopSellingBooks = s.TopSellingBooks(ctx, orders)
	report.GeneratedAt = time.Now()

	if err := s.SaveReportAsJSON(report); err != nil {
		return fmt.Errorf("failed to save report: %v", err)
	}

	return nil
}

func (s *ReportService) TopSellingBooks(ctx context.Context, orders []model.Order) []model.Book {

	bookSales := make(map[int]int)

	for _, order := range orders {
		yesterday := time.Now().AddDate(0, 0, -1)
		if order.CreatedAt.After(yesterday) {

			for _, item := range order.Items {
				bookSales[item.BookID] += item.Quantity
			}
		}
	}

	var sortedSales []model.BookSale

	for bookID, quantity := range bookSales {
		sortedSales = append(sortedSales, model.BookSale{BookID: bookID, Quantity: quantity})
	}

	sort.Slice(sortedSales, func(i, j int) bool {
		return sortedSales[i].Quantity > sortedSales[j].Quantity
	})

	const topN = 3
	topBooks := []model.Book{}

	for i, sale := range sortedSales {
		if i >= topN {
			break
		}

		if book, err := s.BookRepo.GetBook(ctx, sale.BookID); err == nil {
			topBooks = append(topBooks, book)
		}
	}
	return topBooks

}

func (s *ReportService) TotalBooksSold(ctx context.Context, orders []model.Order) int {
	var totalBooks int

	for _, order := range orders {
		yesterday := time.Now().AddDate(0, 0, -1)
		if order.CreatedAt.After(yesterday) {
			for _, item := range order.Items {
				totalBooks += item.Quantity
			}
		}
	}
	return totalBooks

}

func (s *ReportService) TotalOrders(ctx context.Context, orders []model.Order) int {
	var count int

	for _, order := range orders {
		yesterday := time.Now().AddDate(0, 0, -1)
		if order.CreatedAt.After(yesterday) {
			count++
		}
	}
	return count
}

func (s *ReportService) TotalRevenue(ctx context.Context, orders []model.Order) float64 {
	var revenue float64

	for _, order := range orders {
		yesterday := time.Now().AddDate(0, 0, -1)
		if order.CreatedAt.After(yesterday) {
			revenue += order.TotalPrice
		}

	}
	return revenue
}
