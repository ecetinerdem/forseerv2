package main

import (
	"net/http"
)

func (app *application) healthzCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": app.config.version,
	}

	err := app.writeJsonResponse(w, http.StatusOK, data)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
	}
}
