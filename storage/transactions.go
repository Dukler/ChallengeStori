package storage

import (
	"time"

	"github.com/Dukler/ChallengeStori/model"
	"github.com/google/uuid"
)

func (s *PostgresStore) createTransactionTable() error {
	query := `
	create table if not exists transactions (
		id UUID NOT NULL primary key,
		execution_id UUID NOT NULL,
		external_id int NOT NULL,
		value bigint NOT NULL,
		date date NOT NULL,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateTransaction(tx *model.Transaction) (*model.Transaction, error) {
	// Start a new transaction
	dbTx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer dbTx.Rollback()

	query := `
		INSERT INTO transactions
		(id, execution_id, external_id, value, date, created_at) VALUES
		($1, $2, $3, $4, $5, $6)
	`
	id := uuid.New()
	tx.Id = &id
	tx.Created_at = time.Now()

	// Execute the query in the transaction
	_, err = dbTx.Exec(query, tx.Id.String(), tx.ExecutionId, tx.ExternalId, tx.Value, tx.Date, tx.Created_at)
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

	return tx, nil
}
