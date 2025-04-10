package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/satyamkale27/Go-social.git/internal/store"
	"net/http"
	"strconv"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
	User_id int      `json:"user_id"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {

		/*
         note:-

		In the createPostHandler function, the payload of type CreatePostPayload
		gets its values assigned from the HTTP request body provided by the user.
		This is done using the readJSON function, which parses the JSON payload
		from the request body into the payload struct.

		 */

		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content == "" {
		app.badRequestResponse(w, r, fmt.Errorf("content is required"))
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
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
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

	post := getPostFromContext(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.Id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comment = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Posts.Delete(ctx, id); err != nil {
		switch {

		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r) // post received from context not by querying database

	app.store.Posts.

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)

	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postId")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		post, err := app.store.Posts.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)

		/*
									context.WithValue:

									It is used to store a value (post in this case) in the context under
								    a specific key ("post" in this case).
									The new context returned by context.WithValue contains the original
								    context's data along with the new key-value pair.
									This allows the value to be retrieved later using the same key.


									IMP note:-

									In your project, this line is used in the postsContextMiddleware to:

						            Fetch a Post object from the database using the postId from the URL.
									Store the Post object in the context with the key "post".
									Pass the enriched context to the next handler in the chain
							        using next.ServeHTTP(w, r.WithContext(ctx)).
									This allows subsequent handlers or functions (like getPostFromContext)
						            to retrieve the Post object from
							        the context without needing to query the database again.

						 more info:-

						The Post object is stored in the context temporarily to avoid querying
						the database multiple times during the same request

						The context is tied to the lifecycle of a single HTTP request.
						When you store the Post object in the context using context.WithValue,
						it is only accessible while processing that request.
						Once the request is completed, the context (and the data stored in it)
						is discarded.

					   By storing the Post object in the context, you avoid querying the database
					   multiple times for the same data during the same request.
					   For example, if multiple handlers or middleware need access to the same Post object,
					   they can retrieve it from the context instead of querying the database again.


			     	 The data stored in the context is specific to the current request
				      and cannot be shared across different requests.
				      This makes it lightweight and efficient for passing data between middleware and handlers.

		*/

		next.ServeHTTP(w, r.WithContext(ctx))

		/*
			The next.ServeHTTP(w, r.WithContext(ctx)) call passes the modified request
			(with the enriched context) to the next handler in the chain.
			The next parameter is the next handler or middleware in the chain,
			and calling ServeHTTP ensures that the request continues to be processed.
		*/

	})
}

func getPostFromContext(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
