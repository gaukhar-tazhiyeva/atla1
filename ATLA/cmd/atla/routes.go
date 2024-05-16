package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	character1 := r.PathPrefix("/api/v1").Subrouter()

	character1.HandleFunc("/characters", app.getCharacterList).Methods("GET")
	character1.HandleFunc("/characters", app.createCharacterHandler).Methods("POST")
	character1.HandleFunc("/characters/{id:[0-9]+}", app.getCharacterHandler).Methods("GET")
	character1.HandleFunc("/characters/{id:[0-9]+}", app.updateCharacterHandler).Methods("PUT")
	character1.HandleFunc("/characters/{id:[0-9]+}", app.requirePermissions("characters:write", app.deleteCharacterHandler)).Methods("DELETE")

	episode1 := r.PathPrefix("/api/v1").Subrouter()

	episode1.HandleFunc("/episodes", app.getEpisodeList).Methods("GET")
	episode1.HandleFunc("/episodes", app.createEpisodeHandler).Methods("POST")
	episode1.HandleFunc("/episodes/{id:[0-9]+}", app.getEpisodeHandler).Methods("GET")
	episode1.HandleFunc("/episodes/{id:[0-9]+}", app.updateEpisodeHandler).Methods("PUT")
	episode1.HandleFunc("/episodes/{id:[0-9]+}", app.requirePermissions("episodes:write", app.deleteEpisodeHandler)).Methods("DELETE")

	users1 := r.PathPrefix("/api/v1").Subrouter()

	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
