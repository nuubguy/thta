package parser

import (
	"strconv"
	"time"
)

type BankStatement struct {
	UniqueIdentifier string
	Amount           float64
	Date             time.Time
}

type BankStatementParser struct {
	BankStatements []BankStatement
}

func (bp *BankStatementParser) Parse(records [][]string) error {
	for _, record := range records[1:] {
		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return err
		}
		date, err := time.Parse("2006-01-02", record[2])
		if err != nil {
			return err
		}
		bp.BankStatements = append(bp.BankStatements, BankStatement{
			UniqueIdentifier: record[0],
			Amount:           amount,
			Date:             date,
		})
	}
	return nil
}
