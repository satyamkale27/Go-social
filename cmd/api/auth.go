package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/satyamkale27/Go-social.git/internal/store"
	"net/http"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password " validate:"required,min=3,max=72"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.badRequestResponse(w, r, err)
	}

	ctx := r.Context()
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken)) // encrypting the token
	hashToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)

		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)

	}

}
