package main

import (
	"log"
	"thta/parser"
	"thta/service"
)

func main() {
	systemTransactions, err := parser.ParseSystemTransactions()
	if err != nil {
		log.Fatalf("Error parsing system transactions: %v", err)
	}

	bankTransactions, err := parser.ParseBankTransactions()
	if err != nil {
		log.Fatalf("Error parsing bank transactions: %v", err)
	}

	reconciliationService := service.NewReconciliationService()
	reconciliationService.Reconcile(systemTransactions, bankTransactions)
}
