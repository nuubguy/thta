package main

import (
	"fmt"
	"log"
	"thta/model"
	"thta/parser"
)

func main() {
	systemFilePath := "input/system/system_transactions.csv"
	bankFilePath := "input/bank/bank_statements.csv"

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

	// Convert system transactions
	systemTransactions, err := model.ConvertSystemTransactions(systemTransactionParser.Transactions)
	if err != nil {
		log.Fatalf("Error converting system transactions: %v", err)
	}

	// Convert bank transactions
	bankTransactions, err := model.ConvertBankTransactions(bankStatementParser.BankStatements)
	if err != nil {
		log.Fatalf("Error converting bank transactions: %v", err)
	}

	// Print system transactions
	fmt.Println("System Transactions:")
	for _, trx := range systemTransactions {
		fmt.Printf("%+v\n", trx)
	}

	// Print bank transactions
	fmt.Println("\nBank Transactions:")
	for _, trx := range bankTransactions {
		fmt.Printf("%+v\n", trx)
	}
}
