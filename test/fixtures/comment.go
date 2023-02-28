package fixtures

import (
	"go-boilerplate/domain/comment"

	"github.com/brianvoe/gofakeit/v5"
)

// AnyComment returns any comment data
func AnyComment() comment.Comment {
	return comment.Comment{
		ID:           gofakeit.Number(1, 50),
		Description:  gofakeit.HackerPhrase(),
		Type:         comment.Schedule,
		AccountID:    gofakeit.UUID(),
		AdvertiserID: gofakeit.UUID(),
		ListingID:    gofakeit.Numerify("########"),
	}
}
