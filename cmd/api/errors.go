package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJsonError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJsonError(w, http.StatusBadRequest, "not found")
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("resource conflict error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJsonError(w, http.StatusConflict, "the server encountered a conflict")
}
