package parser

import (
	"strconv"
	"time"
)

type Transaction struct {
	TrxID           string
	Amount          int64
	Type            string
	TransactionTime time.Time
	MonthYear       string // Added field for grouping by month
}

type TransactionParser struct {
	Transactions []Transaction
}

func (tp *TransactionParser) Parse(records [][]string) error {
	for _, record := range records[1:] {
		amount, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return err
		}
		transactionTime, err := time.Parse("2006-01-02 15:04:05", record[3])
		if err != nil {
			return err
		}
		monthYear := transactionTime.Format("2006-01")
		tp.Transactions = append(tp.Transactions, Transaction{
			TrxID:           record[0],
			Amount:          amount,
			Type:            record[2],
			TransactionTime: transactionTime,
			MonthYear:       monthYear,
		})
	}
	return nil
}
