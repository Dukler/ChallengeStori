package model

import (
	"time"

	"github.com/google/uuid"
)

type Summary struct {
	ExecutionId   *uuid.UUID     `json:"execution_id"`
	AverageCredit uint64         `json:"average_credit"`
	AverageDebit  uint64         `json:"average_debit"`
	Balance       int64          `json:"balance"`
	TxnsByMonth   map[string]int `json:"monthly_transactions"`
	Created_at    time.Time      `json:"created_at"`
}
