package parser

type CSVParser interface {
	Parse(records [][]string) error
}
