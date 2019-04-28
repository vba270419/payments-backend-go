package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Print("Start Payments Server Application")

	repository, mongoClient, err := InitializeMongoRepository()
	if err != nil {
		os.Exit(1)
	}

	SetPaymentRepository(repository)

	router := ConfigureRouter()

	// TODO move to config
	host := "127.0.0.1"
	port := "8000"

	server := &http.Server{
		Handler:      router,
		Addr:         host + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "graceful server shutdown timeout")
	flag.Parse()

	go func() {
		log.Printf("Starting web server at [%s:%s] ...", host, port)
		log.Fatal(server.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	ShutdownMongoRepository(mongoClient)

	log.Printf("Stopping web server at [%s:%s] ...", host, port)
	_ = server.Shutdown(ctx)
	log.Println("Web server stopped")

	os.Exit(0)
}
