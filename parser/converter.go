package parser

type UnifiedTransaction struct {
	ID         string
	Amount     int64
	Date       string
	Source     string // "System" or "Bank"
	FileSource string // The source file path or name
}

func ConvertSystemTransactions(transactions []Transaction, fileSource string) ([]UnifiedTransaction, error) {
	var unifiedTransactions []UnifiedTransaction
	for _, transaction := range transactions {
		amount := transaction.Amount
		if transaction.Type == "DEBIT" {
			amount = -amount
		}
		unifiedTransactions = append(unifiedTransactions, UnifiedTransaction{
			ID:         transaction.TrxID,
			Amount:     amount,
			Date:       transaction.TransactionTime.Format("2006-01-02"),
			Source:     "System",
			FileSource: fileSource,
		})
	}
	return unifiedTransactions, nil
}

func ConvertBankTransactions(statements []BankStatement, fileSource string) ([]UnifiedTransaction, error) {
	var unifiedTransactions []UnifiedTransaction
	for _, statement := range statements {
		unifiedTransactions = append(unifiedTransactions, UnifiedTransaction{
			ID:         statement.UniqueIdentifier,
			Amount:     statement.Amount,
			Date:       statement.Date.Format("2006-01-02"),
			Source:     "Bank",
			FileSource: fileSource,
		})
	}
	return unifiedTransactions, nil
}
