package main

import (
	"log"
	"thta/parser"
	"thta/service"
	"time"
)

func main() {

	startTime, err := time.Parse("2006-01-02", "2023-01-01")
	if err != nil {
		log.Fatalf("Error parsing start time: %v", err)
	}

	endTime, err := time.Parse("2006-01-02", "2024-12-31")
	if err != nil {
		log.Fatalf("Error parsing end time: %v", err)
	}

	systemTransactions, err := parser.ParseSystemTransactions()
	if err != nil {
		log.Fatalf("Error parsing system transactions: %v", err)
	}

	bankTransactions, err := parser.ParseBankTransactions()
	if err != nil {
		log.Fatalf("Error parsing bank transactions: %v", err)
	}

	reconciliationService := service.NewReconciliationService()
	reconciliationService.Reconcile(systemTransactions, bankTransactions, startTime, endTime)
}
