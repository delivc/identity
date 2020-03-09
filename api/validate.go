package api

import (
	"fmt"
	"net/http"

	"github.com/delivc/identity/models"
	"github.com/gofrs/uuid"
)

type ValidateResponse struct {
	User *models.User `json:"user"`
}

// Validate returns token infos and user informations
// for an internal api endpoint
func (a *API) Validate(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims := getClaims(ctx)
	if claims == nil {
		return badRequestError("Could not read claims")
	}

	userID, err := uuid.FromString(claims.Subject)
	if err != nil {
		return badRequestError("Could not read User ID claim")
	}

	aud := a.requestAud(ctx, r)
	if aud != claims.Audience {
		return badRequestError("Token audience doesn't match request audience")
	}

	user, err := models.FindUserByID(a.db, userID)
	if err != nil {
		if models.IsNotFoundError(err) {
			return notFoundError(err.Error())
		}
		return internalServerError("Database error finding user").WithInternalError(err)
	}

	token := getToken(ctx)
	fmt.Println(token)
	response := ValidateResponse{
		User: user,
	}

	return sendJSON(w, http.StatusOK, response)
}
