package main

import (
	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.Methods("POST").Path("/job").HandlerFunc(createJob)

	return router
}
