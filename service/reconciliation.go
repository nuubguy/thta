package service

import (
	"fmt"
	"sync"
	"thta/model"
	"thta/utils"
)

const (
	DateAmountFormat = "%s_%d"
	NumWorkers       = 4
)

type Task struct {
	SysTrx  model.UnifiedTransaction
	BankMap map[string]model.UnifiedTransaction
}

type ReconciliationService struct{}

func NewReconciliationService() *ReconciliationService {
	return &ReconciliationService{}
}

func (s *ReconciliationService) Reconcile(systemTransactions, bankTransactions []model.UnifiedTransaction) {
	matchedCount, totalDiscrepancy, systemMap, bankMap := s.initializeMaps(systemTransactions, bankTransactions)

	discrepancyChan := make(chan int64)
	taskChan := make(chan model.UnifiedTransaction, len(systemMap))
	var wg sync.WaitGroup

	worker := func(taskChan <-chan model.UnifiedTransaction, discrepancyChan chan<- int64, wg *sync.WaitGroup) {
		defer wg.Done()
		localDiscrepancy := int64(0)

		for sysTrx := range taskChan {
			for bankKey, bankTrx := range bankMap {
				if sysTrx.Date == bankTrx.Date {
					discrepancy := utils.Abs(sysTrx.Amount - bankTrx.Amount)
					localDiscrepancy += discrepancy
					delete(systemMap, fmt.Sprintf(DateAmountFormat, sysTrx.Date, sysTrx.Amount))
					delete(bankMap, bankKey)
					break
				}
			}
		}

		discrepancyChan <- localDiscrepancy
	}

	for i := 0; i < NumWorkers; i++ {
		wg.Add(1)
		go worker(taskChan, discrepancyChan, &wg)
	}

	go func() {
		for _, sysTrx := range systemMap {
			taskChan <- sysTrx
		}
		close(taskChan)
	}()

	go func() {
		wg.Wait()
		close(discrepancyChan)
	}()

	for discrepancy := range discrepancyChan {
		totalDiscrepancy += discrepancy
	}

	s.printResults(matchedCount, totalDiscrepancy, systemMap, bankMap)
}

func (s *ReconciliationService) initializeMaps(systemTransactions, bankTransactions []model.UnifiedTransaction) (int, int64, map[string]model.UnifiedTransaction, map[string]model.UnifiedTransaction) {
	matchedCount := 0
	totalDiscrepancy := int64(0)

	systemMap := make(map[string]model.UnifiedTransaction)
	bankMap := make(map[string]model.UnifiedTransaction)

	for _, trx := range systemTransactions {
		key := fmt.Sprintf(DateAmountFormat, trx.Date, trx.Amount)
		systemMap[key] = trx
	}

	for _, trx := range bankTransactions {
		key := fmt.Sprintf(DateAmountFormat, trx.Date, trx.Amount)
		bankMap[key] = trx
	}

	for key := range systemMap {
		if _, exists := bankMap[key]; exists {
			matchedCount++
			delete(systemMap, key)
			delete(bankMap, key)
		}
	}

	return matchedCount, totalDiscrepancy, systemMap, bankMap
}

func (s *ReconciliationService) printResults(matchedCount int, totalDiscrepancy int64, systemMap, bankMap map[string]model.UnifiedTransaction) {
	fmt.Printf("Total number of matched transactions: %d\n", matchedCount)
	fmt.Printf("Total discrepancies (sum of absolute differences in amount): %d\n", totalDiscrepancy)

	fmt.Printf("\nUnmatched System Transactions: %d\n", len(systemMap))
	for _, sysTrx := range systemMap {
		fmt.Printf("%+v\n", sysTrx)
	}

	unmatchedBankTransactionsByFile := make(map[string][]model.UnifiedTransaction)
	for _, bankTrx := range bankMap {
		unmatchedBankTransactionsByFile[bankTrx.FileSource] = append(unmatchedBankTransactionsByFile[bankTrx.FileSource], bankTrx)
	}

	fmt.Printf("\nUnmatched Bank Transactions: %d\n", len(bankMap))
	for fileSource, transactions := range unmatchedBankTransactionsByFile {
		fmt.Printf("Source File: %s, Unmatched Transactions: %d\n", fileSource, len(transactions))
		for _, bankTrx := range transactions {
			fmt.Printf("%+v\n", bankTrx)
		}
	}
}
