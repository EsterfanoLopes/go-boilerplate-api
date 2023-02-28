package comment

import (
	"database/sql"
	"go-boilerplate/domain/comment"
	"go-boilerplate/repository"
	"go-boilerplate/test/fixtures"
	"testing"
)

var r = &repositoryImpl{}

// Any returns a persisted comment and register its necessary cleanup in the given test
func Any(t *testing.T) comment.Comment {
	cmt := fixtures.AnyComment()
	ID := 0

	repository.Tx(t, func(tx *sql.Tx) {
		id, err := r.Insert(tx, cmt)
		if err != nil {
			t.Errorf("error inserting comment test data %s", err)
		}
		ID = id
	})

	t.Cleanup(func() {
		DeleteTestData(t, ID)
	})
	cmt.ID = ID
	return cmt
}

// DeleteTestData deletes some previous test data
func DeleteTestData(t *testing.T, ID int) {
	repository.Tx(t, func(tx *sql.Tx) {
		err := r.Delete(tx, ID)
		if err != nil {
			t.Errorf("error cleaning up comment test data %s", err)
		}
	})
}

// Comment gets comment test data from database
func Comment(t *testing.T, ID int) comment.Comment {
	cmt := comment.Comment{}
	repository.Tx(t, func(tx *sql.Tx) {
		data, err := r.FindByID(tx, ID)
		if err != nil {
			t.Errorf("error getting comment test data %s", err)
		}
		cmt = data
	})
	return cmt
}
