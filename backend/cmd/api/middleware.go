package main

import (
	"fmt"
	"net/http"
)

// recoverPanic is middleware wrapping the router that ensures we send a 500 Internal Server Error
// in addition to http.Server's panic responses
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
