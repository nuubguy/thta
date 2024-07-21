package parser

import (
	"strconv"
	"thta/constant"
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
		transactionTime, err := time.Parse(constant.DateTimeFormat, record[3])
		if err != nil {
			return err
		}
		monthYear := transactionTime.Format(constant.DateMonthFormat)
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
