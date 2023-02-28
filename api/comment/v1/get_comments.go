package v1

import (
	"go-boilerplate/common/pagination"
	"go-boilerplate/common/response"
	commentFacade "go-boilerplate/facade/comment"
	commentRepository "go-boilerplate/repository/comment"
	"net/http"
)

func parseRequest(r *http.Request) (commentRepository.Query, pagination.Pagination, error) {
	params := r.URL.Query()
	advertiserID := params.Get("advertiserId")
	accountID := params.Get("accountId")
	listingID := params.Get("listingId")

	p, err := pagination.FromRequest(r)
	if err != nil {
		return commentRepository.Query{}, pagination.Pagination{}, err
	}

	return commentRepository.Query{
		AdvertiserID: advertiserID,
		AccountID:    accountID,
		ListingID:    listingID,
	}, p, nil
}

// CommentsGetHandler handle comments get requests
// @Tags Comment
// @Accept json
// @Produce json
// @Param listingId query int false "listingId"
// @Param accountId query int false "accountId"
// @Param advertiserId query int false "advertiserId"
// @Param from query int false "from" Format(int)
// @Param size query int false "size" Format(int)
// @Success 200 {object} pagination.Response{results=[]comment.Comment}
// @Failure 400 {object} response.Error "When some value of the request is invalid"
// @Failure 500 {object} response.Error "When something was wrong when trying to persist data"
// @Router /v1/comment [get]
func CommentsGetHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	q, p, err := parseRequest(r)
	if err != nil {
		response.WriteValidationError(w, err)
		return
	}
	if err := q.Validate(); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	results, count, err := commentFacade.Get().Find(q, p)
	if err != nil {
		response.WriteError(w, r, err, "error finding comments")
		return
	}

	response.Write(w, p.GetResponse(count, results), http.StatusOK)
}
