package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	repository, err := InitializeMongoRepository()
	if err != nil {
		log.Fatal("Failed to initialize mongo repository", err)
		os.Exit(666)
	}

	SetPaymentRepository(repository)
	router := ConfigureRouter()

	log.Print(http.ListenAndServe(":8000", router))
}
