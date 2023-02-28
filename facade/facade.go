package facade

import (
	sql "database/sql"
	"go-boilerplate/repository"
)

var txManager = &TxManagerImpl{}

// TxManager for business logic in facade layer
type TxManager interface {
	// Begin a transaction with database and message buffer
	Begin() (*sql.Tx, error)
	// Resolve given transaction handling message buffer after commit succeeds
	Resolve(*sql.Tx, *error)
}

type TxManagerImpl struct{}

// GetTxManager instance
func GetTxManager() TxManager {
	return txManager
}

func (t *TxManagerImpl) Resolve(tx *sql.Tx, err *error) {
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	} else if *err != nil {
		tx.Rollback()
	} else {
		err := tx.Commit()
		if err != nil {
			panic(err)
		}
	}
}

func (t *TxManagerImpl) Begin() (*sql.Tx, error) {
	tx, err := repository.DB.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// WithTxManager execute given func in a transactional context
func WithTxManager(txm TxManager, fn func(tx *sql.Tx) error) error {
	tx, err := txm.Begin()
	if err != nil {
		return err
	}
	defer txm.Resolve(tx, &err)

	err = fn(tx)
	if err != nil {
		return err
	}

	return nil
}
