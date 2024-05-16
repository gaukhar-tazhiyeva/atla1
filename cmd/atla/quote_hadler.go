package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/justverena/ATLA/pkg/atla/model"
	"github.com/justverena/ATLA/pkg/atla/validator"
)

func (app *application) createQuoteHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Quote string `json:"quote"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	quote := &model.Quote{
		Quote: input.Quote,
	}

	err = app.models.Quotes.Insert(quote)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"quote": quote}, nil)
}

func (app *application) getQuoteList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Quote string
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Quote = app.readStrings(qs, "quote", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		"id", "quote", "created_at", "updated_at",
		"-id", "-quote", "-created_at", "-updated_at",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	quotes, metadata, err := app.models.Quotes.GetAll(input.Quote, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"quotes": quotes, "metadata": metadata}, nil)
}

func (app *application) getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	quote, err := app.models.Quotes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"quote": quote}, nil)
}

func (app *application) updateQuoteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	quote, err := app.models.Quotes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Quote *string `json:"quote"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Quote != nil {
		quote.Quote = *input.Quote
	}
	v := validator.New()

	if model.ValidateQuote(v, quote); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Quotes.Update(quote)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"quote": quote}, nil)
}

func (app *application) deleteQuoteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Quotes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}

func (app *application) getCharacterQuotesList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	characterID, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid character ID")
		return
	}

	character, err := app.models.Characters.Get(characterID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	quotes, err := app.models.Quotes.GetQuotesByCharacterID(characterID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)

	app.writeJSON(w, http.StatusOK, envelope{"quotes": quotes}, nil)

}
