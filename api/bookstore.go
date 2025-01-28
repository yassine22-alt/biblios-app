package main

import (
	"bookstore/api/api/internal/handlers"
	"bookstore/api/api/internal/json"
	"bookstore/api/api/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	bookRepo := json.NewJsonBookStore()
	authorRepo := json.NewJsonAuthorStore()
	customerRepo := json.NewJsonCustomerStore()
	orderRepo := json.NewJsonOrderStore()

	bookService := service.NewBookService(bookRepo, authorRepo)
	authorService := service.NewAuthorService(authorRepo)
	customerService := service.NewCustomerService(customerRepo)
	orderService := service.NewOrderService(orderRepo, customerRepo, bookRepo)
	reportService := service.NewReportService(orderRepo, bookRepo)

	bookHandler := handlers.NewBookHandler(bookService)
	authorHandler := handlers.NewAuthorHandler(authorService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	orderHandler := handlers.NewOrderHandler(orderService)
	reportHandler := handlers.NewReportHandler("./reports")

	//logging
	logFile, err := os.OpenFile("api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	//Middleware for logging http request
	logRequest := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	// mux instead of default serve mux for security purposes !!
	http.Handle("/books", logRequest(http.HandlerFunc(bookHandler.ServeHTTP)))
	http.Handle("/books/{id}", logRequest(http.HandlerFunc(bookHandler.ServeHTTPById)))
	http.Handle("/authors", logRequest(http.HandlerFunc(authorHandler.ServeHTTP)))
	http.Handle("/authors/{id}", logRequest(http.HandlerFunc(authorHandler.ServeHTTPById)))
	http.Handle("/customers", logRequest(http.HandlerFunc(customerHandler.ServeHTTP)))
	http.Handle("/customers/{id}", logRequest(http.HandlerFunc(customerHandler.ServeHTTPById)))
	http.Handle("/orders", logRequest(http.HandlerFunc(orderHandler.ServeHTTP)))
	http.Handle("/orders/{id}", logRequest(http.HandlerFunc(orderHandler.ServeHTTPById)))
	http.Handle("/reports", logRequest(http.HandlerFunc(reportHandler.ServeHTTP)))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go reportService.StartSalesReportGenrator(ctx, logger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		fmt.Println("\nSaving data to files...")

		if err := authorRepo.SaveToFile(); err != nil {
			logger.Printf("Error saving authors: %v\n", err)
		} else if err := customerRepo.SaveToFile(); err != nil {
			logger.Printf("Error saving customers: %v\n", err)
		} else if err := bookRepo.SaveToFile(); err != nil {
			logger.Printf("Error saving books: %v\n", err)
		} else if err := orderRepo.SaveToFile(); err != nil {
			logger.Printf("Error saving orders: %v\n", err)
		}

		fmt.Println("Data saved successfully")
		os.Exit(0)
	}()

	err1 := http.ListenAndServe(":8080", nil)
	if err1 != nil {
		logger.Println("Error serving:", err1)
	}

}
