/*!!!!!!!!!!!!!!!!!! HAS CRUD !!!!!!!!!!!!!!!!!!*/

package main

import (
	"errors"
	"log"
	"net/http"

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
		// ID        int
		// Gender    string
		// Status    string
		// Nation    string
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	// Use our helpers to extract the title and nutrition value range query string values, falling back to the
	// defaults of an empty string and an empty slice, respectively, if they are not provided
	// by the client.
	input.Name = app.readStrings(qs, "name", "")
	input.Age = app.readInt(qs, "age", 0, v)
	// input.Gender = app.readStrings(qs, "gender", "")
	// input.Status = app.readStrings(qs, "status", "")
	// input.Nation = app.readStrings(qs, "nation", "")
	// Ge the page and page_size query string value as integers. Notice that we set the default
	// page value to 1 and default page_size to 20, and that we pass the validator instance
	// as the final argument.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply an ascending sort on menu ID).
	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	// Add the supported sort value for this endpoint to the sort safelist.
	// name of the column in the database.
	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "name", "age",
		// descending sort values
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

// func (app *application) readJSON(_ http.ResponseWriter, r *http.Request, dst interface{}) error {
// 	dec := json.NewDecoder(r.Body)
// 	dec.DisallowUnknownFields()

// 	err := dec.Decode(dst)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
