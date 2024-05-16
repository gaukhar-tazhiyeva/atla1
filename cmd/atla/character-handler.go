/*!!!!!!!!!!!!!!!!!! HAS CRUD !!!!!!!!!!!!!!!!!!*/

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

func (app *application) createCharacterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
		Status string `json:"status"`
		Nation string `json:"nation"`
		// CreatedAt string `json:"createdAt"`
		// UpdatedAt string `json:"updatedAt"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	character := &model.Character{
		ID:     input.ID,
		Name:   input.Name,
		Age:    input.Age,
		Gender: input.Gender,
		Status: input.Status,
		Nation: input.Nation,
		// CreatedAt: input.CreatedAt,
		// UpdatedAt: input.UpdatedAt,
	}

	err = app.models.Characters.Insert(character)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"character": character}, nil)
}

func (app *application) getCharacterList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		Age  int
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readStrings(qs, "name", "")
	input.Age = app.readInt(qs, "age", 0, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		"id", "name", "age",
		"-id", "-name", "-age",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	characters, metadata, err := app.models.Characters.GetAll(input.Name, input.Age, input.Age, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"characters": characters, "metadata": metadata}, nil)
}

func (app *application) getCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	character, err := app.models.Characters.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)
}

func (app *application) updateCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	character, err := app.models.Characters.Get(id)
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
		ID     *int    `json:"id"`
		Name   *string `json:"name"`
		Age    *int    `json:"age"`
		Gender *string `json:"gender"`
		Status *string `json:"status"`
		Nation *string `json:"nation"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		character.Name = *input.Name
	}

	if input.Age != nil {
		character.Age = *input.Age
	}
	if input.Gender != nil {
		character.Gender = *input.Gender
	}
	if input.Status != nil {
		character.Status = *input.Status
	}
	if input.Nation != nil {
		character.Nation = *input.Nation
	}
	v := validator.New()

	if model.ValidateCharacter(v, character); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Characters.Update(character)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)
}

func (app *application) deleteCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Characters.Delete(id)
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

func (app *application) getEpisodeCharacters(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	episodeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid episode ID")
		return
	}

	episode, err := app.models.Episodes.Get(episodeID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	characters, err := app.models.Characters.GetByEpisode(episodeID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"episode": episode}, nil)

	app.writeJSON(w, http.StatusOK, envelope{"characters": characters}, nil)
}

func (app *application) getQuoteCharacter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	quoteID, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid quote ID")
		return
	}

	quote, err := app.models.Quotes.Get(quoteID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	character, err := app.models.Characters.GetByQuote(quoteID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"quote": quote}, nil)

	app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)
}
