package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal  error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")

}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())

}

func (app *application) conflictResponce(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusConflict, err.Error())

}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusNotFound, "not found")

}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")

}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("www-authenticate", `Basic realm="restricted, charset=UTF-8"`)

	/*
		The line of code that caused the browser's login prompt box is:
		w.Header().Set("www-authenticate", `Basic realm="restricted, charset=UTF-8"`)
		This is located in the unauthorizedBasicErrorResponse function in the file cmd/api/errors.go.
		It sets the WWW-Authenticate header, which instructs the browser to display the login prompt
		when a 401 Unauthorized response is sent.
	*/
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")

}
