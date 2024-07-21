package parser

import (
	"strconv"
	"time"
)

type Transaction struct {
	TrxID           string
	Amount          float64
	Type            string
	TransactionTime time.Time
}

type TransactionParser struct {
	Transactions []Transaction
}

func (tp *TransactionParser) Parse(records [][]string) error {
	for _, record := range records[1:] {
		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return err
		}
		transactionTime, err := time.Parse("2006-01-02 15:04:05", record[3])
		if err != nil {
			return err
		}
		tp.Transactions = append(tp.Transactions, Transaction{
			TrxID:           record[0],
			Amount:          amount,
			Type:            record[2],
			TransactionTime: transactionTime,
		})
	}
	return nil
}
