package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chefgoldbloom/pnctool/internal/data"
	"github.com/chefgoldbloom/pnctool/internal/validator"
)

// Add a createCameraHandler for the "POST /v1/cameras" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createCameraHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name       string `json:"name"`
		MacAddress string `json:"mac_address"`
		SiteName   string `json:"site_name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	camera := &data.Camera{
		Name:       input.Name,
		MacAddress: input.MacAddress,
		SiteName:   input.SiteName,
	}
	v := validator.New()
	if data.ValidateCamera(v, camera); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showCameraHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	camera := data.Camera{
		ID:         id,
		CreatedAt:  time.Now(),
		Name:       "This-Is-TestSite Test Camera",
		MacAddress: "ACCC8593034234",
		SiteName:   "This-Is-TestSite",
		Username:   "root",
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"camera": camera}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
