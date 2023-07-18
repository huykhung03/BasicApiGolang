package sqlc

import (
	"context"
	"database/sql"
	"fmt"
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

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// * Second argument in store.db.BeginTx() is isolation level in database
	// * store.db.BeginTx() returns a transaction object (tx object) or error
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rollbackErr)
		}
		return err
	}

	return tx.Commit()
}
