package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"thta/model"
	"thta/parser"
)

func reconcileTransactions(systemTransactions []model.UnifiedTransaction, bankTransactions []model.UnifiedTransaction) {
	matchedCount := 0
	totalDiscrepancy := int64(0)

	// Maps to track unmatched transactions
	systemMap := make(map[string]model.UnifiedTransaction)
	bankMap := make(map[string]model.UnifiedTransaction)

	// Populate system transactions map
	for _, trx := range systemTransactions {
		key := fmt.Sprintf("%s_%d", trx.Date, trx.Amount)
		systemMap[key] = trx
	}

	// Populate bank transactions map
	for _, trx := range bankTransactions {
		key := fmt.Sprintf("%s_%d", trx.Date, trx.Amount)
		bankMap[key] = trx
	}

	// Identify matched transactions
	for key, _ := range systemMap {
		if _, exists := bankMap[key]; exists {
			matchedCount++
			delete(systemMap, key)
			delete(bankMap, key)
		}
	}

	// Calculate discrepancies and remove processed entries
	for date, sysTrx := range systemMap {
		for bankKey, bankTrx := range bankMap {
			if sysTrx.Date == bankTrx.Date {
				discrepancy := abs(sysTrx.Amount - bankTrx.Amount)
				totalDiscrepancy += discrepancy
				delete(systemMap, date)
				delete(bankMap, bankKey)
				break
			}
		}
	}

	// Print the number of matched transactions
	fmt.Printf("Total discrepancies: %d\n", totalDiscrepancy)
	fmt.Printf("Total number of matched transactions: %d\n", matchedCount)
	// Identify unmatched system transactions
	fmt.Printf("\nUnmatched System Transactions: %d\n", len(systemMap))

	for _, sysTrx := range systemMap {
		fmt.Printf("%+v\n", sysTrx)
	}

	// Group unmatched bank transactions by FileSource
	unmatchedBankTransactionsByFile := make(map[string][]model.UnifiedTransaction)
	for _, bankTrx := range bankMap {
		unmatchedBankTransactionsByFile[bankTrx.FileSource] = append(unmatchedBankTransactionsByFile[bankTrx.FileSource], bankTrx)
	}

	// Print the number of unmatched bank transactions grouped by FileSource
	fmt.Printf("\nUnmatched Bank Transactions: %d\n", len(bankMap))
	for fileSource, transactions := range unmatchedBankTransactionsByFile {
		fmt.Printf("Source File: %s, Unmatched Transactions: %d\n", fileSource, len(transactions))
		for _, bankTrx := range transactions {
			fmt.Printf("%+v\n", bankTrx)
		}
	}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	systemFilePath := "input/system/system_transactions.csv"
	bankDirPath := "input/bank/"

	systemTransactionParser := &parser.TransactionParser{}
	err := parser.ParseCSV(systemFilePath, systemTransactionParser)
	if err != nil {
		log.Fatalf("Error parsing system transactions CSV: %v", err)
	}

	var bankTransactions []model.UnifiedTransaction

	// Convert system transactions
	systemTransactions, err := model.ConvertSystemTransactions(systemTransactionParser.Transactions, systemFilePath)
	if err != nil {
		log.Fatalf("Error converting system transactions: %v", err)
	}

	err = filepath.Walk(bankDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			bankStatementParser := &parser.BankStatementParser{}
			err = parser.ParseCSV(path, bankStatementParser)
			if err != nil {
				log.Fatalf("Error parsing bank statement CSV: %v", err)
			}
			transactions, err := model.ConvertBankTransactions(bankStatementParser.BankStatements, path)
			if err != nil {
				log.Fatalf("Error converting bank transactions: %v", err)
			}
			bankTransactions = append(bankTransactions, transactions...)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error reading bank statement files: %v", err)
	}

	reconcileTransactions(systemTransactions, bankTransactions)
}
