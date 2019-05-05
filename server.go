package main

import (
	"context"
	"flag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	serverHost    string = "server_host"
	serverPort    string = "server_port"
	serverTimeout string = "server_timeout"

	mongoDbHost    string = "mongodb_host"
	mongoDbPort    string = "mongodb_port"
	mongoDbTimeout string = "mongodb_timeout"
)

func main() {
	log.Print("Start Payments Server Application")

	initializeEnvironmentProperties()

	repository, mongoClient := InitializeMongoRepository()

	SetPaymentRepository(repository)

	router := configureRouter()

	host := viper.GetString(serverHost)
	port := viper.GetString(serverPort)
	timeout := time.Duration(viper.GetInt(serverTimeout)) * time.Second

	server := &http.Server{
		Handler:      router,
		Addr:         host + ":" + port,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
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

func initializeEnvironmentProperties() {
	log.Print("Checking environment properties...")

	configurationFile := flag.String("conf", "./config/server.json", "Path to configuration file")
	flag.Parse()

	viper.SetConfigFile(*configurationFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if !viper.IsSet(serverHost) {
		log.Fatal("Server host property is not configured")
	}

	if !viper.IsSet(serverPort) {
		log.Fatal("Server port property is not configured")
	}

	if !viper.IsSet(mongoDbTimeout) {
		log.Fatal("MongoDB timeout property is not configured")
	}

	if !viper.IsSet(mongoDbHost) {
		log.Fatal("MongoDB host property is not configured")
	}

	if !viper.IsSet(mongoDbPort) {
		log.Fatal("MongoDB port property is not configured")
	}

	if !viper.IsSet(mongoDbTimeout) {
		log.Fatal("MongoDB timeout property is not configured")
	}

	log.Print("Environment properties - OK")
}
