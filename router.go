package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

const (
	methodPost   string = "POST"
	methodPut    string = "PUT"
	methodDelete string = "DELETE"
	methodGet    string = "GET"

	createPaymentPath  string = "/v1/payments/create"
	updatePaymentPath  string = "/v1/payments/update"
	deletePaymentPath  string = "/v1/payments/delete/{id}"
	getPaymentPath     string = "/v1/payments/get/{id}"
	getAllPaymentsPath string = "/v1/payments/all"
)

type route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

var routes []route

func initializeRoutes() {
	addRoute(route{createPaymentPath, methodPost, createPaymentEndpoint})
	addRoute(route{updatePaymentPath, methodPut, updatePaymentEndpoint})
	addRoute(route{deletePaymentPath, methodDelete, deletePaymentEndpoint})
	addRoute(route{getPaymentPath, methodGet, getPaymentEndpoint})
	addRoute(route{getAllPaymentsPath, methodGet, getAllPaymentsEndpoint})
}

func addRoute(route route) {
	routes = append(routes, route)
}

func configureRouter() (router *mux.Router) {
	log.Print("Initializing router...")

	initializeRoutes()

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

func preparePaymentURL(path string, paymentID string) string {
	return strings.ReplaceAll(path, "{id}", paymentID)
}

func prepareFullPaymentURL(host string, path string, paymentID string) string {
	return fmt.Sprintf("http://%s%s", host, preparePaymentURL(path, paymentID))
}
