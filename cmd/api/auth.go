package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/satyamkale27/Go-social.git/internal/mailer"
	"github.com/satyamkale27/Go-social.git/internal/store"
	"net/http"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
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

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
		// it is used to just return token to user
	}

	acticationUrl := fmt.Sprintf("%s/confirm/%s", app.config.frontendUrl, plainToken)

	isProdEnv := app.config.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: acticationUrl,
	}

	err = app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending mail", "error", err)
		// rollback user creation if email (fails SAGA pattern)
		if err := app.store.Users.Delete(ctx, user.Id); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
			app.internalServerError(w, r, err)
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)

	}

}
