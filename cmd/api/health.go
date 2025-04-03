package main

import (
	"log"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {

		log.Printf(err.Error())

	}

}
