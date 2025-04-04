package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/satyamkale27/Go-social.git/internal/store"
	"net/http"
	"strconv"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	User_id int      `json:"user_id"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())

		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}
	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	/*
			The two writeJSONError calls handle different error scenarios:


			The first writeJSONError handles errors that occur during the creation of the post
		    in the database. If app.store.Posts.Create(ctx, post) fails, it returns an error,
		    and the function responds with an internal server error status and the error message.


			The second writeJSONError handles errors that occur while writing the JSON response
		    back to the client. If writeJSON(w, http.StatusCreated, post) fails, it returns an error,
		    and the function responds with an internal server error status and the error message.

	*/

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := r.Context()
	post, err := app.store.Posts.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
