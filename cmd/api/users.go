package main

import (
	"fmt"
	"net/http"
)

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "register handler")
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login handler")
}
