package main

import (
	"log"
	"thta/parser"
)

func main() {
	systemFilePath := "system_transactions.csv"
	bankFilePath := "bank_statements.csv"

	systemTransactionParser := &parser.TransactionParser{}
	err := parser.ParseCSV(systemFilePath, systemTransactionParser)
	if err != nil {
		log.Fatalf("Error parsing system transactions CSV: %v", err)
	}

	bankStatementParser := &parser.BankStatementParser{}
	err = parser.ParseCSV(bankFilePath, bankStatementParser)
	if err != nil {
		log.Fatalf("Error parsing bank statements CSV: %v", err)
	}
}
