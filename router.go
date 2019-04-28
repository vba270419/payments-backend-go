package main

import (
	"github.com/gorilla/mux"
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
	InitializeRoutes()

	router = mux.NewRouter()
	for _, route := range routes {
		router.HandleFunc(route.Path, route.Handler).Methods(route.Method)
	}

	router.StrictSlash(true)
	return router
}
