package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/satyamkale27/Go-social.git/internal/store"
	"net/http"
	"strconv"
)

type userKey string

const userCtx = "user"

type CreateUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getPostFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUser struct {
	userID int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {

	followerUser := getPostFromContext(r)

	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context().Value
	app.store.Users.follow(ctx, followerUser.Id, payload.userID)
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	user := getPostFromContext(r)
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, userId)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponse(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUserFromContext(r *http.Request) store.User {
	user, _ := r.Context().Value(userCtx).(store.User)
	return user
}
