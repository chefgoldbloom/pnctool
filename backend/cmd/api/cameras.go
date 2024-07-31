package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chefgoldbloom/pnctool/backend/internal/data"
	"github.com/chefgoldbloom/pnctool/backend/internal/validator"
)

// Add a createCameraHandler for the "POST /v1/cameras" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createCameraHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name       string `json:"name"`
		MacAddress string `json:"mac_address"`
		SiteName   string `json:"site_name"`
		ModelNo    string `json:"model_no"`
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
		ModelNo:    input.ModelNo,
	}
	v := validator.New()
	if data.ValidateCamera(v, camera); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Cameras.Insert(camera)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// Make a Location header to let the client know resource's url
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/cameras/%d", camera.ID))

	// Write JSON response with a 201 Created status code, the camera data
	// in the response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"camera": camera}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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

	camera, err := app.models.Cameras.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"camera": camera}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCameraHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	camera, err := app.models.Cameras.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name       *string `json:"name"`
		MacAddress *string `json:"mac_address"`
		SiteName   *string `json:"site_name"`
		ModelNo    *string `json:"model_no"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	if input.Name != nil {
		camera.Name = *input.Name
	}
	if input.MacAddress != nil {
		camera.MacAddress = *input.MacAddress
	}
	if input.SiteName != nil {
		camera.SiteName = *input.SiteName

	}
	if input.ModelNo != nil {
		camera.ModelNo = *input.ModelNo
	}

	v := validator.New()
	if data.ValidateCamera(v, camera); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Cameras.Update(camera)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"camera": camera}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCameraHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Cameras.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "camera successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
