package main

import (
	"API-Mutant/routing"
	"net/http"
	"log"
	"os"
)

func main(){
	router := routing.NewRouter()

	server := http.ListenAndServe(":"+os.Args[1], router)

	log.Fatal(server)
}