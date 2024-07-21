package parser

import (
	"encoding/csv"
	"os"
)

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
