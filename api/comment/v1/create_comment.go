// Package v1 holds comment v1 api handlers
package v1

import (
	"encoding/json"
	"go-boilerplate/common/response"
	"go-boilerplate/domain/comment"
	commentFacade "go-boilerplate/facade/comment"
	"net/http"
)

// CommentPostHandler handle comment post requests
// @Tags Comment
// @Accept json
// @Produce json
// @Param comment body comment.Comment true "payload"
// @Success 201 {object} response.Success
// @Failure 400 {object} response.Error "When some value of the request is invalid"
// @Failure 500 {object} response.Error "When something was wrong when trying to persist data"
// @Router /v1/comment [post]
func CommentPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body := comment.Comment{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.WriteUnprocessableEntity(w, err)
		return
	}
	if err := body.Validate(); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	ID, err := commentFacade.Get().Insert(body)
	if err != nil {
		response.WriteError(w, r, err, "error inserting comment")
		return
	}

	response.Write(w, response.Success{
		ID: ID,
	}, http.StatusCreated)
}
