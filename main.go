package main

import (
	"API-Mutant/routing"
	"net/http"
	"log"
)

func main(){
	router := routing.NewRouter()

	server := http.ListenAndServe(":3000", router)

	log.Fatal(server)
}