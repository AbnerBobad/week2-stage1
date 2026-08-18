package main

import (
	"net/http"

	"week2-stage1/internal/data"
	"week2-stage1/internal/validator"
)

func (app *application) createConsumersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.Name != "", "name", "must be provided")
	v.Check(input.Email != "", "email", "must be provided")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	consumer := &data.Consumer{
		CN_Name:  input.Name,
		CN_Email: input.Email,
	}

	err = app.models.Consumers.Insert(consumer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"consumer": consumer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
