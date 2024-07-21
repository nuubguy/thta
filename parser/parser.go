package parser

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
)

const (
	systemFilePath = "input/system/system_transactions.csv"
	bankDirPath    = "input/bank/"
)

func ParseSystemTransactions() ([]UnifiedTransaction, error) {
	systemTransactionParser := &TransactionParser{}
	err := ParseCSV(systemFilePath, systemTransactionParser)
	if err != nil {
		return nil, err
	}

	return ConvertSystemTransactions(systemTransactionParser.Transactions, systemFilePath)
}

func ParseBankTransactions() ([]UnifiedTransaction, error) {
	var bankFiles []string
	err := filepath.Walk(bankDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			bankFiles = append(bankFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var bankTransactions []UnifiedTransaction
	for _, path := range bankFiles {
		bankStatementParser := &BankStatementParser{}
		err := ParseCSV(path, bankStatementParser)
		if err != nil {
			log.Fatalf("Error parsing bank statement CSV: %v", err)
		}
		transactions, err := ConvertBankTransactions(bankStatementParser.BankStatements, path)
		if err != nil {
			log.Fatalf("Error converting bank transactions: %v", err)
		}
		bankTransactions = append(bankTransactions, transactions...)
	}

	return bankTransactions, nil
}

func ParseCSV(filePath string, parser CSVParser) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	return parser.Parse(records)
}
