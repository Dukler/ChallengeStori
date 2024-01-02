package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dukler/ChallengeStori/model"
)

func (s *PostgresStore) createSummaryTable() error {
	query := "CREATE EXTENSION IF NOT EXISTS hstore;"
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	query = `
	create table if not exists summary (
		execution_id UUID NOT NULL,
		average_credit bigint NOT NULL,
		average_debit bigint NOT NULL,
		balance bigint NOT NULL,
		monthly_transactions hstore NOT NULL,
		created_at timestamp
	)`

	_, err = s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateSummary(sum *model.Summary) (*model.Summary, error) {
	// Start a new transaction
	dbTx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer dbTx.Rollback()
	sum.Created_at = time.Now()
	query := `
		INSERT INTO summary 
		(execution_id, average_credit, average_debit, balance, monthly_transactions, created_at) VALUES
		($1, $2, $3, $4, $5, $6)
	`

	// Build the hstore string directly from the map
	hstoreString := ""
	for key, value := range sum.TxnsByMonth {
		hstoreString += fmt.Sprintf(`"%s"=>"%d", `, key, value)
	}
	hstoreString = strings.TrimSuffix(hstoreString, ", ")

	// Execute the query in the transaction
	_, err = dbTx.Exec(query, sum.ExecutionId, sum.AverageCredit, sum.AverageDebit, sum.Balance, hstoreString, sum.Created_at)
	if err != nil {
		// If there's an error, rollback the transaction
		dbTx.Rollback()
		return nil, err
	}

	// If everything goes well, commit the transaction
	err = dbTx.Commit()
	if err != nil {
		return nil, err
	}

	return sum, nil
}
