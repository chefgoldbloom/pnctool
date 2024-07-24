package main

import (
	"fmt"
	"net/http"
)

// Add a createCameraHandler for the "POST /v1/cameras" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createCameraHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a camera")
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showCameraHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {

		http.NotFound(w, r)
	}

	fmt.Fprintf(w, "show camera with id: %d\n", id)
}
