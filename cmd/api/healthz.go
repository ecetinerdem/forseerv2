package main

import (
	"log"
	"net/http"
)

func (app *application) healthzCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": app.config.version,
	}

	err := writeJson(w, http.StatusOK, data)
	if err != nil {
		log.Println(err.Error())
	}
}
