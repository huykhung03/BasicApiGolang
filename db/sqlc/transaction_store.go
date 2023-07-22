package sqlc

import (
	"context"
	"database/sql"
)

// Store provides all functions to excute db queries and transaction
type Store interface {
	Querier
	PurchaseTransaction(ctx context.Context, arg PurchaseTransactionPagrams) (PurchaseTransactionResult, error)
}

// SQLStore provides all functions to excute SQL queries and transaction
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}
