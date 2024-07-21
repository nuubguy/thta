package model

import (
	"thta/parser"
)

type UnifiedTransaction struct {
	ID     string
	Amount int64
	Date   string
	Source string // "System" or "Bank"
}

func ConvertSystemTransactions(transactions []parser.Transaction) ([]UnifiedTransaction, error) {
	var unifiedTransactions []UnifiedTransaction
	for _, transaction := range transactions {
		amount := transaction.Amount
		if transaction.Type == "DEBIT" {
			amount = -amount
		}
		unifiedTransactions = append(unifiedTransactions, UnifiedTransaction{
			ID:     transaction.TrxID,
			Amount: amount,
			Date:   transaction.TransactionTime.Format("2006-01-02"),
			Source: "System",
		})
	}
	return unifiedTransactions, nil
}

func ConvertBankTransactions(statements []parser.BankStatement) ([]UnifiedTransaction, error) {
	var unifiedTransactions []UnifiedTransaction
	for _, statement := range statements {
		unifiedTransactions = append(unifiedTransactions, UnifiedTransaction{
			ID:     statement.UniqueIdentifier,
			Amount: statement.Amount,
			Date:   statement.Date.Format("2006-01-02"),
			Source: "Bank",
		})
	}
	return unifiedTransactions, nil
}
