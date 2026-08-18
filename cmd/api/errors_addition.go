package main

import "net/http"

// Helper method for sending an error response to the client with status code.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := app.writeJSON(w, status, envelope{"error": message}, nil)
	if err != nil {
		app.logger.Error(err.Error())
		w.WriteHeader(500)
	}
}

// serverErrorResponse - logs details of the given error and sends a 500 error
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	message := "Server Encountered an error, could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// Helper method for sending a 404 Not Found response to the client.
// Uses the errorResponse helper method to send a JSON response with the appropriate status code and message.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// 400 - Bad request response helper method
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// Failed validation response helper method
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
