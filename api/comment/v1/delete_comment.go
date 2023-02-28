package v1

import (
	"go-boilerplate/common/response"
	commentFacade "go-boilerplate/facade/comment"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CommentDeleteHandler handle comment delete requests
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path int true "id" Format(int)
// @Success 204 {string} string "Success"
// @Failure 400 {object} response.Error "When some value of the request is invalid"
// @Failure 500 {object} response.Error "When something was wrong when trying to persist data"
// @Router /v1/comment/{id} [delete]
func CommentDeleteHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		response.WriteValidationError(w, err)
		return
	}

	err = commentFacade.Get().Delete(ID)
	if err != nil {
		response.WriteError(w, r, err, "error deleting comment")
		return
	}

	response.Write(w, nil, http.StatusNoContent)
}
