package comment_test

import (
	"database/sql"
	"go-boilerplate/common/pagination"
	"go-boilerplate/domain/comment"
	"go-boilerplate/facade"
	commentFacade "go-boilerplate/facade/comment"
	commentRepository "go-boilerplate/repository/comment"
	"go-boilerplate/test/fixtures"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"
)

var (
	txManagerMock = &facade.MockTxManager{}
	commentsMock  = &commentRepository.MockRepository{}
	f             = commentFacade.Facade{
		TxManager: txManagerMock,
		Comments:  commentsMock,
	}
	verifyAllMocks = func(t *testing.T) {
		txManagerMock.AssertExpectations(t)
		commentsMock.AssertExpectations(t)
	}
)

func TestInsert(t *testing.T) {
	cmt := fixtures.AnyComment()
	testCases := []struct {
		name           string
		comment        comment.Comment
		configureMocks func()
	}{
		{
			name:    "comment inserted successfully",
			comment: cmt,
			configureMocks: func() {
				txManagerMock.On("Begin").Return(nil, nil, nil).Once()
				txManagerMock.On("Resolve", mock.Anything, mock.Anything, mock.Anything).Once()

				commentsMock.On("Insert", mock.Anything, cmt).Return(cmt.ID, nil).Once()
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.configureMocks()

			_, err := f.Insert(tc.comment)
			if err != nil {
				t.Errorf("error inserting comment %s", err)
			}

			verifyAllMocks(t)
		})
	}
}

func TestUpdate(t *testing.T) {
	cmt := fixtures.AnyComment()
	testCases := []struct {
		name           string
		comment        comment.Comment
		configureMocks func()
	}{
		{
			name:    "comment updated successfully",
			comment: cmt,
			configureMocks: func() {
				txManagerMock.On("Begin").Return(nil, nil, nil).Once()
				txManagerMock.On("Resolve", mock.Anything, mock.Anything, mock.Anything).Once()

				commentsMock.On("Update", mock.Anything, cmt).Return(nil).Once()
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.configureMocks()

			err := f.Update(tc.comment)
			if err != nil {
				t.Errorf("error updating comment %s", err)
			}

			verifyAllMocks(t)
		})
	}
}

func TestFindByID(t *testing.T) {
	cmt := fixtures.AnyComment()
	testCases := []struct {
		name           string
		ID             int
		expected       comment.Comment
		configureMocks func()
	}{
		{
			name:     "found comment by id",
			ID:       cmt.ID,
			expected: cmt,
			configureMocks: func() {
				commentsMock.On("FindByID", (*sql.Tx)(nil), cmt.ID).Return(cmt, nil).Once()
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.configureMocks()

			result, err := f.FindByID(tc.ID)
			if err != nil {
				t.Errorf("unexpected error finding comment by id %s", err)
				return
			}
			if diff := cmp.Diff(result, tc.expected, cmpopts.IgnoreFields(comment.Comment{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("unexpected comment %s", diff)
			}

			verifyAllMocks(t)
		})
	}
}

func TestFind(t *testing.T) {
	cmt := fixtures.AnyComment()

	q := commentRepository.Query{
		AdvertiserID: cmt.AdvertiserID,
		AccountID:    cmt.AccountID,
		ListingID:    cmt.ListingID,
	}
	p, _ := pagination.New(0, 10)

	testCases := []struct {
		name           string
		query          commentRepository.Query
		pagination     pagination.Pagination
		configureMocks func()
		expected       []comment.Comment
		expectedCount  int
	}{
		{
			name:       "some comments found",
			query:      q,
			pagination: p,
			configureMocks: func() {
				commentsMock.On("Find", (*sql.Tx)(nil), q, p).Return([]comment.Comment{
					cmt,
				}, nil).Once()
				commentsMock.On("Count", (*sql.Tx)(nil), q).Return(1, nil).Once()
			},
			expected: []comment.Comment{
				cmt,
			},
			expectedCount: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.configureMocks()

			results, count, err := f.Find(tc.query, tc.pagination)
			if err != nil {
				t.Errorf("error finding comments %s", err)
				return
			}

			if diff := cmp.Diff(results, tc.expected, cmpopts.IgnoreFields(comment.Comment{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("unexpected comments %s", diff)
				return
			}

			if count != tc.expectedCount {
				t.Errorf("unexpected comments count %d", count)
				return
			}

			verifyAllMocks(t)
		})
	}
}

func TestDelete(t *testing.T) {
	ID := gofakeit.Number(1, 10)
	testCases := []struct {
		name           string
		ID             int
		configureMocks func()
	}{
		{
			name: "comment deleted successfully",
			ID:   ID,
			configureMocks: func() {
				txManagerMock.On("Begin").Return(nil, nil, nil).Once()
				txManagerMock.On("Resolve", mock.Anything, mock.Anything, mock.Anything).Once()

				commentsMock.On("Delete", mock.Anything, ID).Return(nil).Once()
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.configureMocks()

			err := f.Delete(tc.ID)
			if err != nil {
				t.Errorf("error deleting comment %s", err)
			}

			verifyAllMocks(t)
		})
	}
}
