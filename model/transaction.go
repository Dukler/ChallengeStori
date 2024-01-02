package model

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id          *uuid.UUID `json:"id"`
	ExecutionId *uuid.UUID `json:"execution_id"`
	ExternalId  int        `json:"external_id"`
	Value       int64      `json:"value"`
	Date        time.Time  `json:"date"`
	Created_at  time.Time  `json:"created_at"`
}
