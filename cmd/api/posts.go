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

						the readJSON function is used in the createPostHandler
			            to parse the JSON payload from the HTTP request body into the CreatePostPayload
			            struct. This allows the application to extract and validate the data provided by
			            the user in the request.
			            The readJSON function decodes the JSON payload into the struct and ensures that the
			            data matches the expected structure. If the JSON is invalid or contains unexpected
			            fields, it returns an error, which is then handled by the badRequestResponse function.

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

	user := getUserFromContext(r) // get the user that is currently authenticated

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.Id,
	}
	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
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

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
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

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r) // post received from context not by querying database

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	/*

				IMP NOTE:-
				In `UpdatePostPayload`, the fields are defined as pointers (e.g., `*string`) so that you can tell the difference between:

				1. **Field not provided by the user**:
				   - If the user doesn't include a field in their request, its value will be `nil`.
				   - This means the user doesn't want to update that field.

				2. **Field provided but empty**:
				   - If the user includes a field but leaves it empty (e.g., `"title": ""`), the pointer will not be `nil`. Instead, it will point to an empty string (`""`).
				   - This means the user wants to update the field and set it to an empty value.

				### Why is this important?
				When updating a post, you need to know whether:
				- The user wants to **skip updating a field** (leave it as it is).
				- The user wants to **update a field and set it to an empty value**.

				Using pointers helps you handle this distinction.

				### Example:
				#### Request 1: User skips the `title` field
				```json
				{
				  "content": "Updated content"
				}
				```
				- `payload.Title` will be `nil` (not provided).
				- You won't update the `title` field in the database.

				#### Request 2: User provides an empty `title`
				```json
				{
				  "title": "",
				  "content": "Updated content"
				}
				```
				- `payload.Title` will point to `""` (empty string).
				- You will update the `title` field in the database and set it to an empty value.

				This flexibility is why pointers are used in `UpdatePostPayload`.




			************************************* what if pointer is not used  **********************************************************


				If you don't use pointers for fields in the `UpdatePostPayload` struct, you won't be able to differentiate between a field that is **not provided**
		        in the request and a field that is **provided with an empty value**. Here's what would happen:

				### Without Pointers:
				```go
				type UpdatePostPayload struct {
				    Title   string `json:"title" validate:"omitempty,max=100"`
				    Content string `json:"content" validate:"omitempty,max=1000"`
				}
				```

				1. **Default Values**:
				   - If a field is not provided in the JSON request, it will be assigned the default value for its type:
				     - For `string`, the default value is an empty string (`""`).
				   - This makes it impossible to know whether the user explicitly set the field to an empty value or simply omitted it.

				2. **Behavior**:
				   - If the user sends a request like this:
				     ```json
				     {
				       "title": "New Title"
				     }
				     ```
				     - `payload.Content` will be an empty string (`""`), even though the user didn't provide it.
				     - This would overwrite the `Content` field in the database with an empty value, which is likely not the intended behavior.

				3. **No Partial Updates**:
				   - Without pointers, you cannot perform partial updates because you cannot distinguish between "no update" and "update to an empty value."

				### With Pointers:
				Using pointers allows you to check for `nil` to determine if a field was provided in the request. If a field is `nil`, you can skip updating it.

				### Example Comparison:
				#### Without Pointers:
				```go
				if payload.Content != "" {
				    post.Content = payload.Content
				}
				```
				- This will always update `post.Content` to an empty string if the user doesn't provide the `content` field.

				#### With Pointers:
				```go
				if payload.Content != nil {
				    post.Content = *payload.Content
				}
				```
				- This will only update `post.Content` if the user explicitly provides the `content` field in the request.

				### Conclusion:
				Using pointers is essential for handling partial updates correctly. Without them, you lose the ability to distinguish between "field not provided" and "field provided with an empty value," which can lead to unintended overwrites.


	*/

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
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
