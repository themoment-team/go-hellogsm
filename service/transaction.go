package service

import (
	"database/sql"
)

type TransactionManager interface {
	BeginTx() (tx interface{}, err error)
	Commit(tx interface{}) error
	Rollback(tx interface{}) error
}

// DefaultTxManager default TransactionManager implementation
type DefaultTxManager struct {
	db *sql.DB
}

// NewTransactionManager create a TransactionManager instance
func NewTransactionManager(db *sql.DB) TransactionManager {
	return &DefaultTxManager{
		db: db,
	}
}

// BeginTx begin a transaction
func (tm *DefaultTxManager) BeginTx() (interface{}, error) {
	tx, err := tm.db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// Commit a transaction
func (tm *DefaultTxManager) Commit(tx interface{}) error {
	tx1 := tx.(*sql.Tx)
	err := tx1.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Rollback a transaction
func (tm *DefaultTxManager) Rollback(tx interface{}) error {
	tx1 := tx.(*sql.Tx)
	err := tx1.Rollback()
	if err != nil {
		return err
	}
	return nil
}
