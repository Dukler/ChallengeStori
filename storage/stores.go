package storage

import (
	"database/sql"
	"log/slog"

	"github.com/Dukler/ChallengeStori/model"
)

type PostgresStore struct {
	db *sql.DB
	l  *slog.Logger
}

type transactionStore interface {
	CreateTransaction(*model.Transaction) (*model.Transaction, error)
}

type summaryStore interface {
	CreateSummary(*model.Summary) (*model.Summary, error)
}
