package storage

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	dbname   = "storiDB"
	user     = "postgres"
	password = "postgrespw"
)

type Storage interface {
	transactionStore
	summaryStore
}

func NewPostgresStore(maxRetries int, retryInterval time.Duration, logger *slog.Logger) (*PostgresStore, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var db *sql.DB
	var err error

	for i := 1; i <= maxRetries; i++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			logger.Warn("Error opening database: %v", err)
		} else if err = db.Ping(); err == nil {
			logger.Info("Connected to the database")
			break
		}

		logger.Warn(fmt.Sprintf("Error connecting to the database (attempt %d): %s", i, err.Error()))
		db.Close()
		time.Sleep(retryInterval * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database after %d retries: %v", maxRetries, err)
	}
	maxConnections := 40
	db.SetMaxOpenConns(maxConnections)
	db.SetMaxIdleConns(maxConnections)

	return &PostgresStore{
		db: db,
		l:  logger,
	}, nil
}

func (s *PostgresStore) Init() error {
	if err := s.createTransactionTable(); err != nil {
		return err
	}
	if err := s.createSummaryTable(); err != nil {
		return err
	}
	s.l.Info("Migration finished.")
	return nil
}
