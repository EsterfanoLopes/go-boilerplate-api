// Package comment holds comment business logic
package comment

import (
	"database/sql"
	"go-boilerplate/common/pagination"
	"go-boilerplate/domain/comment"
	"go-boilerplate/facade"
	commentRepository "go-boilerplate/repository/comment"
)

var (
	instance = &Facade{
		TxManager: facade.GetTxManager(),
		Comments:  commentRepository.Get(),
	}
)

type Facade struct {
	TxManager facade.TxManager
	Comments  commentRepository.Repository
}

func Get() *Facade {
	return instance
}

// Insert a comment
func (f *Facade) Insert(cmt comment.Comment) (ID int, err error) {
	err = facade.WithTxManager(f.TxManager, func(tx *sql.Tx) error {
		var err error
		ID, err = f.Comments.Insert(tx, cmt)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

// Update a comment
func (f *Facade) Update(cmt comment.Comment) error {
	return facade.WithTxManager(f.TxManager, func(tx *sql.Tx) error {
		return f.Comments.Update(tx, cmt)
	})
}

// FindByID a comment
func (f *Facade) FindByID(ID int) (comment.Comment, error) {
	return f.Comments.FindByID(nil, ID)
}

// Find and count comments given a query
func (f *Facade) Find(q commentRepository.Query, p pagination.Pagination) ([]comment.Comment, int, error) {
	results, err := f.Comments.Find(nil, q, p)
	if err != nil {
		return nil, 0, err
	}

	count, err := f.Comments.Count(nil, q)
	if err != nil {
		return nil, 0, err
	}

	return results, count, err
}

// Delete a comment
func (f *Facade) Delete(ID int) (err error) {
	return facade.WithTxManager(f.TxManager, func(tx *sql.Tx) error {
		return f.Comments.Delete(tx, ID)
	})
}
