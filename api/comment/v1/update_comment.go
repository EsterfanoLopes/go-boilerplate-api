package v1

import (
	"encoding/json"
	"go-boilerplate/common/response"
	"go-boilerplate/domain/comment"
	commentFacade "go-boilerplate/facade/comment"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CommentPutHandler handle comment put requests
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path int true "id" Format(int)
// @Param comment body comment.Comment true "payload"
// @Success 200 {object} response.Success
// @Failure 400 {object} response.Error "When some value of the request is invalid"
// @Failure 500 {object} response.Error "When something was wrong when trying to persist data"
// @Router /v1/comment/{id} [put]
func CommentPutHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		response.WriteValidationError(w, err)
		return
	}

	body := comment.Comment{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.WriteUnprocessableEntity(w, err)
		return
	}
	if err := body.Validate(); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	body.ID = ID
	err = commentFacade.Get().Update(body)
	if err != nil {
		response.WriteError(w, r, err, "error updating comment")
		return
	}

	response.Write(w, response.Success{
		ID: ID,
	}, http.StatusOK)
}
