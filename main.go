package main

import (
	"API-Mutant/routing"
	"net/http"
	"log"
)

func main(){
	router := routing.NewRouter()

	server := http.ListenAndServe(":8080", router)

	log.Fatal(server)
}