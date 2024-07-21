package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"thta/model"
	"thta/parser"
	"thta/service"
)

func parseAndConvertBankFiles(files []string, bankTransactionsChan chan<- []model.UnifiedTransaction, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, path := range files {
		bankStatementParser := &parser.BankStatementParser{}
		err := parser.ParseCSV(path, bankStatementParser)
		if err != nil {
			log.Fatalf("Error parsing bank statement CSV: %v", err)
		}
		transactions, err := model.ConvertBankTransactions(bankStatementParser.BankStatements, path)
		if err != nil {
			log.Fatalf("Error converting bank transactions: %v", err)
		}
		bankTransactionsChan <- transactions
	}
}

func main() {
	systemFilePath := "input/system/system_transactions.csv"
	bankDirPath := "input/bank/"

	systemTransactionParser := &parser.TransactionParser{}
	err := parser.ParseCSV(systemFilePath, systemTransactionParser)
	if err != nil {
		log.Fatalf("Error parsing system transactions CSV: %v", err)
	}

	systemTransactions, err := model.ConvertSystemTransactions(systemTransactionParser.Transactions, systemFilePath)
	if err != nil {
		log.Fatalf("Error converting system transactions: %v", err)
	}

	var bankFiles []string
	err = filepath.Walk(bankDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			bankFiles = append(bankFiles, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error reading bank statement files: %v", err)
	}

	bankTransactionsChan := make(chan []model.UnifiedTransaction)
	var wg sync.WaitGroup

	numWorkers := 4
	fileBatches := len(bankFiles) / numWorkers
	for i := 0; i < numWorkers; i++ {
		start := i * fileBatches
		end := (i + 1) * fileBatches
		if i == numWorkers-1 {
			end = len(bankFiles)
		}
		wg.Add(1)
		go parseAndConvertBankFiles(bankFiles[start:end], bankTransactionsChan, &wg)
	}

	go func() {
		wg.Wait()
		close(bankTransactionsChan)
	}()

	var bankTransactions []model.UnifiedTransaction
	for transactions := range bankTransactionsChan {
		bankTransactions = append(bankTransactions, transactions...)
	}

	reconciliationService := service.NewReconciliationService()
	reconciliationService.Reconcile(systemTransactions, bankTransactions)
}
