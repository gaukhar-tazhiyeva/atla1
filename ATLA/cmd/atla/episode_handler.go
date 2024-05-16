package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/justverena/ATLA/pkg/atla/model"
	"github.com/justverena/ATLA/pkg/atla/validator"
)

func (app *application) createEpisodeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Air_Date string `json:"air_date"`
		// CreatedAt string `json:"createdAt"`
		// UpdatedAt string `json:"updatedAt"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	episode := &model.Episode{
		ID:       input.ID,
		Title:    input.Title,
		Air_Date: input.Air_Date,
		// CreatedAt: input.CreatedAt,
		// UpdatedAt: input.UpdatedAt,
	}

	err = app.models.Episodes.Insert(episode)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"episodes": episode}, nil)
}

func (app *application) getEpisodeList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string
		Air_Date string
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	// Use our helpers to extract the title and nutrition value range query string values, falling back to the
	// defaults of an empty string and an empty slice, respectively, if they are not provided
	// by the client.
	input.Title = app.readStrings(qs, "title", "")
	input.Air_Date = app.readStrings(qs, "air_date", "")
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
		"id", "title", "air_date",
		// descending sort values
		"-id", "-title", "-air_date",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	episodes, metadata, err := app.models.Episodes.GetAll(input.Title, input.Air_Date, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"episodes": episodes, "metadata": metadata}, nil)
}

func (app *application) getEpisodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	episode, err := app.models.Episodes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"episodes": episode}, nil)
}

func (app *application) updateEpisodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	episode, err := app.models.Episodes.Get(id)
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
		ID       *int    `json:"id"`
		Title    *string `json:"title"`
		Air_Date *string `json:"air_date"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		episode.Title = *input.Title
	}

	if input.Air_Date != nil {
		episode.Air_Date = *input.Air_Date
	}
	v := validator.New()

	if model.ValidateEpisode(v, episode); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Episodes.Update(episode)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"episodes": episode}, nil)
}

func (app *application) deleteEpisodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Episodes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "episode deleted successfully"}, nil)
}
