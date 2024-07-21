# Transaction Reconciliation Service

This application is designed to reconcile transactions between internal system transactions and external bank statements. It identifies unmatched and discrepant transactions within a specified date range.

## Setup

### Directory Structure

Place your transaction CSV files in the following directories:

- **System Transactions**: `input/system/system_transactions.csv`
- **Bank Transactions**: `input/bank/`

### Input File Format

#### System Transactions

- File: `input/system/system_transactions.csv`
- Format: CSV with the following columns:
    - `trxID` (string): Unique identifier for the transaction
    - `amount` (int64): Transaction amount
    - `type` (string): Transaction type (DEBIT or CREDIT)
    - `transactionTime` (string): Date and time of the transaction (format: YYYY-MM-DD HH:MM:SS)

#### Bank Statements

- Files: Place multiple bank statement CSV files in the `input/bank/` directory.
- Format: CSV with the following columns:
    - `unique_identifier` (string): Unique identifier for the transaction in the bank statement
    - `amount` (int64): Transaction amount (can be negative for debits)
    - `date` (string): Date of the transaction (format: YYYY-MM-DD)

### Example Files

#### `input/system/system_transactions.csv`

```csv
trxID,amount,type,transactionTime
1,1000,CREDIT,2023-01-05 10:00:00
2,-500,DEBIT,2023-02-15 14:30:00
```

#### `input/bank/bank_statement_1.csv`

```csv
unique_identifier,amount,date
abc123,1000,2023-01-05
def456,-500,2023-02-15
```