package main

import (
	"github.com/gorilla/mux"
)

// Gorilla Mux is the router.  The package is listed above.
// The items listed in "{}" are variables that can be set in the route
func NewRouter(app *AppContext) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/person/export", use(AppHandler{app, ExportCSV})).Methods("GET")
	router.Handle(`/person/{person:[\w+\-\s+.]+}`, use(AppHandler{app, ListPerson})).Name("person").Methods("GET")
	router.Handle("/person", use(AppHandler{app, AddPerson})).Methods("POST")
	router.Handle("/person", use(AppHandler{app, UpdatePerson})).Methods("PUT")
	router.Handle("/person", use(AppHandler{app, RemovePerson})).Methods("DELETE")
	router.Handle("/person/upload", use(AppHandler{app, ImportCSV})).Methods("POST")
	router.Handle("/people", use(AppHandler{app, ListAllEntries})).Methods("GET")

	return router
}
