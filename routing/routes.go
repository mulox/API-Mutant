package routing

import (
	"github.com/gorilla/mux"
	"API-Mutant/controllers"
)


func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	
	r.HandleFunc("/mutant", controllers.IsMutant).Methods("POST")
	r.HandleFunc("/stats", controllers.GetStats).Methods("GET")

	return r
}