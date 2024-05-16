package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	character1 := r.PathPrefix("/api/v1").Subrouter()

	character1.HandleFunc("/characters/{id:[0-9]+}/episode", app.getCharacterEpisode).Methods("GET")
	character1.HandleFunc("/characters/{id:[0-9]+}/quotes", app.getCharacterQuotesList).Methods("GET")

	character1.HandleFunc("/characters", app.getCharacterList).Methods("GET")
	character1.HandleFunc("/characters", app.createCharacterHandler).Methods("POST")
	character1.HandleFunc("/characters/{id:[0-9]+}", app.getCharacterHandler).Methods("GET")
	character1.HandleFunc("/characters/{id:[0-9]+}", app.updateCharacterHandler).Methods("PUT")
	character1.HandleFunc("/characters/{id:[0-9]+}", app.requirePermissions("characters:write", app.deleteCharacterHandler)).Methods("DELETE")

	episode1 := r.PathPrefix("/api/v1").Subrouter()

	episode1.HandleFunc("/episodes/{id:[0-9]+}/characters", app.getEpisodeCharacters).Methods("GET")

	episode1.HandleFunc("/episodes", app.getEpisodeList).Methods("GET")
	episode1.HandleFunc("/episodes", app.createEpisodeHandler).Methods("POST")
	episode1.HandleFunc("/episodes/{id:[0-9]+}", app.getEpisodeHandler).Methods("GET")
	episode1.HandleFunc("/episodes/{id:[0-9]+}", app.updateEpisodeHandler).Methods("PUT")
	episode1.HandleFunc("/episodes/{id:[0-9]+}", app.requirePermissions("episodes:write", app.deleteEpisodeHandler)).Methods("DELETE")

	quote1 := r.PathPrefix("/api/v1").Subrouter()

	episode1.HandleFunc("/quotes/{id:[0-9]+}/character", app.getQuoteCharacter).Methods("GET")

	quote1.HandleFunc("/quotes", app.getQuoteList).Methods("GET")
	quote1.HandleFunc("/quotes", app.createQuoteHandler).Methods("POST")
	quote1.HandleFunc("/quotes/{id:[0-9]+}", app.getQuoteHandler).Methods("GET")
	quote1.HandleFunc("/quotes/{id:[0-9]+}", app.updateQuoteHandler).Methods("PUT")
	quote1.HandleFunc("/quotes/{id:[0-9]+}", app.requirePermissions("quotes:write", app.deleteQuoteHandler)).Methods("DELETE")

	users1 := r.PathPrefix("/api/v1").Subrouter()

	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	return app.authenticate(r)
}
