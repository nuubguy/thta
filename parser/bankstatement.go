package parser

import (
	"strconv"
	"time"
)

type BankStatement struct {
	UniqueIdentifier string
	Amount           int64
	Date             time.Time
	MonthYear        string // Added field for grouping by month
}

type BankStatementParser struct {
	BankStatements []BankStatement
}

func (bp *BankStatementParser) Parse(records [][]string) error {
	for _, record := range records[1:] {
		amount, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return err
		}
		date, err := time.Parse("2006-01-02", record[2])
		if err != nil {
			return err
		}
		monthYear := date.Format("2006-01")
		bp.BankStatements = append(bp.BankStatements, BankStatement{
			UniqueIdentifier: record[0],
			Amount:           amount,
			Date:             date,
			MonthYear:        monthYear,
		})
	}
	return nil
}
