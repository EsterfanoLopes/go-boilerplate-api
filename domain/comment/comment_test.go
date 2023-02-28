package comment_test

import (
	"go-boilerplate/domain/comment"
	"go-boilerplate/test"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		name          string
		comment       comment.Comment
		expectedError string
	}{
		{
			name: "valid comment",
			comment: comment.Comment{
				Type:         comment.Lead,
				Description:  gofakeit.HackerPhrase(),
				AdvertiserID: gofakeit.UUID(),
				AccountID:    gofakeit.UUID(),
				ListingID:    gofakeit.Numerify("##########"),
			},
		},
		{
			name: "missing type",
			comment: comment.Comment{
				Description:  gofakeit.HackerPhrase(),
				AdvertiserID: gofakeit.UUID(),
				AccountID:    gofakeit.UUID(),
				ListingID:    gofakeit.Numerify("##########"),
			},
			expectedError: "type: cannot be blank.",
		},
		{
			name: "missing description",
			comment: comment.Comment{
				Type:         comment.Lead,
				AdvertiserID: gofakeit.UUID(),
				AccountID:    gofakeit.UUID(),
				ListingID:    gofakeit.Numerify("##########"),
			},
			expectedError: "description: cannot be blank.",
		},
		{
			name: "missing advertiser id",
			comment: comment.Comment{
				Type:        comment.Lead,
				Description: gofakeit.HackerPhrase(),
				AccountID:   gofakeit.UUID(),
				ListingID:   gofakeit.Numerify("##########"),
			},
			expectedError: "advertiserId: cannot be blank.",
		},
		{
			name: "invalid advertiser id",
			comment: comment.Comment{
				Type:         comment.Lead,
				Description:  gofakeit.HackerPhrase(),
				AdvertiserID: gofakeit.Fruit(),
				AccountID:    gofakeit.UUID(),
				ListingID:    gofakeit.Numerify("##########"),
			},
			expectedError: "advertiserId: must be a valid UUID.",
		},
		{
			name: "missing account id",
			comment: comment.Comment{
				Type:         comment.Lead,
				Description:  gofakeit.HackerPhrase(),
				AdvertiserID: gofakeit.UUID(),
				ListingID:    gofakeit.Numerify("##########"),
			},
			expectedError: "accountId: cannot be blank.",
		},
		{
			name: "invalid account id",
			comment: comment.Comment{
				Type:         comment.Lead,
				Description:  gofakeit.HackerPhrase(),
				AdvertiserID: gofakeit.UUID(),
				AccountID:    gofakeit.Animal(),
				ListingID:    gofakeit.Numerify("##########"),
			},
			expectedError: "accountId: must be a valid UUID.",
		},
		{
			name: "missing listing id",
			comment: comment.Comment{
				Type:         comment.Lead,
				Description:  gofakeit.HackerPhrase(),
				AdvertiserID: gofakeit.UUID(),
				AccountID:    gofakeit.UUID(),
			},
			expectedError: "listingId: cannot be blank.",
		},
		{
			name: "invalid listing id",
			comment: comment.Comment{
				Type:         comment.Lead,
				Description:  gofakeit.HackerPhrase(),
				AdvertiserID: gofakeit.UUID(),
				AccountID:    gofakeit.UUID(),
				ListingID:    gofakeit.Adverb(),
			},
			expectedError: "listingId: must contain digits only.",
		},
		{
			name: "invalid owner",
			comment: comment.Comment{
				Type:         comment.Lead,
				Description:  gofakeit.HackerPhrase(),
				AdvertiserID: gofakeit.UUID(),
				AccountID:    gofakeit.UUID(),
				ListingID:    gofakeit.Numerify("##########"),
			},
			expectedError: "owner: (accountId: cannot be blank; email: cannot be blank; name: cannot be blank.).",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.comment.Validate()
			test.AssertError(t, err, tc.expectedError)
		})
	}
}
