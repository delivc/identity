package api

import (
	"net/http"

	"github.com/delivc/identity/models"
	"github.com/gofrs/uuid"
)

// ValidateResponse params for the validate response
type ValidateResponse struct {
	Expires int64        `json:"exp"`
	User    *models.User `json:"user"`
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

	response := ValidateResponse{
		Expires: claims.ExpiresAt,
		User:    user,
	}

	return sendJSON(w, http.StatusOK, response)
}
