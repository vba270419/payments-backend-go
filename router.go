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

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

var routes []Route

func InitializeRoutes() {
	AddRoute(Route{createPaymentPath, methodPost, CreatePaymentHandler})
	AddRoute(Route{updatePaymentPath, methodPut, UpdatePaymentHandler})
	AddRoute(Route{deletePaymentPath, methodDelete, DeletePaymentHandler})
	AddRoute(Route{getPaymentPath, methodGet, GetPaymentHandler})
	AddRoute(Route{getAllPaymentsPath, methodGet, GetAllPaymentsHandler})
}

func AddRoute(route Route) {
	routes = append(routes, route)
}

func ConfigureRouter() (router *mux.Router) {
	log.Print("Initializing router...")

	InitializeRoutes()

	router = mux.NewRouter()
	router.Use(LoggingMiddleware)

	for _, route := range routes {
		router.HandleFunc(route.Path, route.Handler).Methods(route.Method)
	}

	router.StrictSlash(true)
	log.Print("Router - OK")
	return router
}

func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("[%s] %s", request.Method, request.RequestURI)
		handler.ServeHTTP(writer, request)
	})
}

func PreparePaymentURL(path string, paymentId string) string {
	return strings.ReplaceAll(path, "{id}", paymentId)
}

func PrepareFullPaymentURL(host string, path string, paymentId string) string {
	return fmt.Sprintf("http://%s%s", host, PreparePaymentURL(path, paymentId))
}
