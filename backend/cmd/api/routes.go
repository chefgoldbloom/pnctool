package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Use Custom error handlers
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Endpoints for cameras
	router.HandlerFunc(http.MethodGet, "/v1/cameras", app.listCamerasHandler)
	router.HandlerFunc(http.MethodPost, "/v1/cameras", app.createCameraHandler)
	router.HandlerFunc(http.MethodGet, "/v1/cameras/:id", app.showCameraHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/cameras/:id", app.updateCameraHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/cameras/:id", app.deleteCameraHandler)

	return app.recoverPanic(router)
}
