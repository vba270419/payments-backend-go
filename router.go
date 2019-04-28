package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

var routes []Route

func InitializeRoutes() {
	AddRoute(Route{"/payments/v1/payment", "POST", CreatePaymentHandler})
	AddRoute(Route{"/payments/v1/payment", "PUT", UpdatePaymentHandler})
	AddRoute(Route{"/payments/v1/payment/{id}", "DELETE", DeletePaymentHandler})
	AddRoute(Route{"/payments/v1/payment/{id}", "GET", GetPaymentHandler})
	AddRoute(Route{"/payments/v1/payments", "GET", GetAllPaymentsHandler})
}

func AddRoute(route Route) {
	routes = append(routes, route)
}

func ConfigureRouter() (router *mux.Router) {
	log.Print("Initializing router...")

	InitializeRoutes()

	router = mux.NewRouter()
	router.Use(loggingMiddleware)

	for _, route := range routes {
		router.HandleFunc(route.Path, route.Handler).Methods(route.Method)
	}

	router.StrictSlash(true)
	log.Print("Router - OK")
	return router
}

func loggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("[%s] %s", request.Method, request.RequestURI)
		handler.ServeHTTP(writer, request)
	})
}
