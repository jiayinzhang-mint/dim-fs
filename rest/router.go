package rest

import "github.com/gorilla/mux"

// InitImageAPI config
func InitImageAPI(r *mux.Router) {
	publicRouter := r.PathPrefix("/image").Subrouter()
	publicRouter.HandleFunc("", viewImage).Queries("path", "{path}").Methods("GET")
	publicRouter.HandleFunc("", uploadImage).Methods("POST")
}
