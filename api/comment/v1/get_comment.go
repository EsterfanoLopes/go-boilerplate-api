package v1

import (
	"go-boilerplate/common/response"
	commentFacade "go-boilerplate/facade/comment"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CommentGetHandler handle comment get requests
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path int true "id" Format(int)
// @Success 200 {object} comment.Comment
// @Failure 400 {object} response.Error "When some value of the request is invalid"
// @Failure 500 {object} response.Error "When something was wrong when trying to persist data"
// @Router /v1/comment/{id} [get]
func CommentGetHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		response.WriteValidationError(w, err)
		return
	}

	result, err := commentFacade.Get().FindByID(ID)
	if err != nil {
		response.WriteError(w, r, err, "error finding comment by id")
		return
	}

	response.Write(w, result, http.StatusOK)
}
