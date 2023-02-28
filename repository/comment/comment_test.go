package comment_test

import (
	"database/sql"
	"go-boilerplate/common/pagination"
	"go-boilerplate/domain/comment"
	"go-boilerplate/repository"
	commentRepository "go-boilerplate/repository/comment"
	"go-boilerplate/test"
	"go-boilerplate/test/fixtures"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var impl = commentRepository.Get()

func TestMain(m *testing.M) {
	err := repository.Setup()
	if err != nil {
		os.Exit(-1)
	}
	os.Exit(m.Run())
}

func TestValidateQuery(t *testing.T) {
	testCases := []struct {
		name        string
		query       commentRepository.Query
		expectedErr string
	}{
		{
			name: "valid query",
			query: commentRepository.Query{
				AccountID:    gofakeit.UUID(),
				AdvertiserID: gofakeit.UUID(),
			},
		},
		{
			name: "missing account id",
			query: commentRepository.Query{
				AdvertiserID: gofakeit.UUID(),
			},
			expectedErr: "AccountID: cannot be blank.",
		},
		{
			name: "invalid account id",
			query: commentRepository.Query{
				AccountID:    gofakeit.BeerHop(),
				AdvertiserID: gofakeit.UUID(),
			},
			expectedErr: "AccountID: must be a valid UUID.",
		},
		{
			name: "missing advertiser id",
			query: commentRepository.Query{
				AccountID: gofakeit.UUID(),
			},
			expectedErr: "AdvertiserID: cannot be blank.",
		},
		{
			name: "invalid advertiser id",
			query: commentRepository.Query{
				AccountID:    gofakeit.UUID(),
				AdvertiserID: gofakeit.Fruit(),
			},
			expectedErr: "AdvertiserID: must be a valid UUID.",
		},
		{
			name: "invalid listing id",
			query: commentRepository.Query{
				AccountID:    gofakeit.UUID(),
				AdvertiserID: gofakeit.UUID(),
				ListingID:    gofakeit.AppAuthor(),
			},
			expectedErr: "ListingID: must contain digits only.",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.query.Validate()
			test.AssertError(t, err, tc.expectedErr)
		})
	}
}

func TestInsert(t *testing.T) {
	nextID := repository.GetNextID(t, "comment")
	cmt := fixtures.AnyComment()

	testCases := []struct {
		name     string
		comment  comment.Comment
		expected int
	}{
		{
			name:     "comment inserted successfully",
			comment:  cmt,
			expected: nextID,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository.Tx(t, func(tx *sql.Tx) {
				ID, err := impl.Insert(tx, tc.comment)
				if err != nil {
					t.Errorf("unexpected error inserting comment %s", err)
					return
				}

				if ID != tc.expected {
					t.Errorf("unexpected next comment id %d", ID)
					return
				}
			})
		})
	}
}

func TestUpdate(t *testing.T) {
	cmt := commentRepository.Any(t)
	cmt.Updated = true
	cmt.Description = gofakeit.Phrase()

	testCases := []struct {
		name     string
		comment  comment.Comment
		expected comment.Comment
	}{
		{
			name:     "comment updated successfully",
			comment:  cmt,
			expected: cmt,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository.Tx(t, func(tx *sql.Tx) {
				err := impl.Update(tx, tc.expected)
				if err != nil {
					t.Errorf("unexpected error updating comment %s", err)
					return
				}
			})

			result := commentRepository.Comment(t, tc.expected.ID)
			if diff := cmp.Diff(result, tc.expected, cmpopts.IgnoreFields(comment.Comment{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("unexpected comment result %s", diff)
				return
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	cmt := commentRepository.Any(t)
	testCases := []struct {
		name     string
		ID       int
		expected comment.Comment
	}{
		{
			name:     "found comment by id",
			ID:       cmt.ID,
			expected: cmt,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository.Tx(t, func(tx *sql.Tx) {
				l, err := impl.FindByID(tx, tc.ID)
				if err != nil {
					t.Errorf("unexpected error finding comment %s", err)
					return
				}
				if diff := cmp.Diff(l, tc.expected, cmpopts.IgnoreFields(comment.Comment{}, "CreatedAt", "UpdatedAt")); diff != "" {
					t.Errorf("unexpected comment %s", diff)
					return
				}
			})
		})
	}
}

func TestFind(t *testing.T) {
	cmt := commentRepository.Any(t)
	p, _ := pagination.New(0, 30)

	testCases := []struct {
		name     string
		q        commentRepository.Query
		p        pagination.Pagination
		expected []comment.Comment
	}{
		{
			name: "found comments by query",
			q: commentRepository.Query{
				AccountID:    cmt.AccountID,
				AdvertiserID: cmt.AdvertiserID,
				ListingID:    cmt.ListingID,
			},
			p: p,
			expected: []comment.Comment{
				cmt,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository.Tx(t, func(tx *sql.Tx) {
				result, err := impl.Find(tx, tc.q, tc.p)
				if err != nil {
					t.Errorf("unexpected error finding comments %s", err)
					return
				}

				if diff := cmp.Diff(result, tc.expected, cmpopts.IgnoreFields(comment.Comment{}, "CreatedAt", "UpdatedAt")); diff != "" {
					t.Errorf("unexpected comments result %s", diff)
					return
				}
			})
		})
	}
}

func TestCount(t *testing.T) {
	cmt := commentRepository.Any(t)

	testCases := []struct {
		name     string
		q        commentRepository.Query
		expected int
	}{
		{
			name: "count 1 comment by query",
			q: commentRepository.Query{
				AccountID:    cmt.AccountID,
				AdvertiserID: cmt.AdvertiserID,
				ListingID:    cmt.ListingID,
			},
			expected: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository.Tx(t, func(tx *sql.Tx) {
				result, err := impl.Count(tx, tc.q)
				if err != nil {
					t.Errorf("unexpected error counting comments %s", err)
					return
				}

				if result != tc.expected {
					t.Errorf("unexpected comments count result %d", result)
					return
				}
			})
		})
	}
}

func TestDelete(t *testing.T) {
	cmt := commentRepository.Any(t)

	testCases := []struct {
		name string
		ID   int
	}{
		{
			name: "comment deleted successfully",
			ID:   cmt.ID,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repository.Tx(t, func(tx *sql.Tx) {
				err := impl.Delete(tx, tc.ID)
				if err != nil {
					t.Errorf("unexpected error deleting comment %s", err)
					return
				}
			})
		})
	}
}
