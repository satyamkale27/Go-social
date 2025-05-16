package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/satyamkale27/Go-social.git/internal/store"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {

			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)

		userid, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, userid)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
		}

		ctx = context.WithValue(ctx, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {

				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicErrorResponse(w, r, err)
				return
			}

			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			/*
					note:-
					In Go, a struct is a blueprint, and when you assign values to it, you create an instance of that struct.

					Explanation:
					Struct Definition: A struct is just a type definition (blueprint) and does not hold any data until an instance is created.


					type basicConfig struct {
					    user string
					    pass string
					}
					Instance of Struct: When you assign values to a struct (e.g., via a variable or field), you create an instance of it. For example:

					basic := basicConfig{user: "admin", pass: "password"}
					Accessing Fields: When you access app.config.auth.basic.user, you are accessing
				    the user field of an instance of the basicConfig struct.
				    This instance is part of the authConfig struct, which is part of the config
				    struct, and so on.


					How to Know if You Are Pointing to a Struct or an Instance:
					If you are accessing fields (e.g., app.config.auth.basic.user), you are working with an instance of the struct.
					If you are referring to the struct type itself (e.g., basicConfig), you are referring to the struct definition (blueprint).

			*/

			cred := strings.SplitN(string(decoded), ":", 2)
			if len(cred) != 2 || cred[0] != username || cred[1] != pass {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r) // it will be the authenticated user
		post := getPostFromContext(r)

		// if it is the users post, is this the owner of the post
		if post.UserID == user.Id {
			next.ServeHTTP(w, r)
			return
		}

		// role precedence check

		allowed, err := app.checkRoleprecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		if !allowed {
			app.forbiddenResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRoleprecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {

	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}
