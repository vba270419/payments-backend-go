package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var paymentRepository PaymentRepository

func setPaymentRepository(repository PaymentRepository) {
	paymentRepository = repository
}

func writeHeaderLocation(writer http.ResponseWriter, request *http.Request, paymentID string) {
	location := prepareFullPaymentURL(request.Host, getPaymentPath, paymentID)
	writer.Header().Set("Location", location)
}

func prepareSuccessHeader(writer http.ResponseWriter, statusCode int) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(statusCode)
}

func prepareFailureHeader(writer http.ResponseWriter, request *http.Request, err error) {
	log.Printf("Request [%s] %s with processed error `%s`", request.Method, request.RequestURI, err.Error())

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")

	switch err.(type) {
	case *PersistenceError:
		writer.WriteHeader(http.StatusInternalServerError)
		return
	case *PaymentNotFoundError:
		writer.WriteHeader(http.StatusNotFound)
		return
	case *PaymentVersionConflictError:
		writer.WriteHeader(http.StatusConflict)
		return
	case *InvalidPaymentError:
		writer.WriteHeader(http.StatusBadRequest)
		return
	default:
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}

// At the moment payment validation has very simple rules: OrganisationID is a required field and payment ID should be not empty for an update call
// here we can add more complex validation rules if needed
func decodeAndValidatePayment(request *http.Request, create bool) (payment Payment, err error) {
	err = json.NewDecoder(request.Body).Decode(&payment)
	if err != nil {
		return payment, err
	}

	if len(payment.OrganisationID) == 0 || (!create && len(payment.ID) == 0) {
		return payment, &InvalidPaymentError{payment}
	}

	return payment, nil
}

func createPaymentEndpoint(writer http.ResponseWriter, request *http.Request) {
	var payment, err = decodeAndValidatePayment(request, true)
	if err != nil {
		prepareFailureHeader(writer, request, err)
		return
	}

	newUUID, _ := uuid.NewUUID()
	payment.ID = newUUID.String()
	payment.Version = 1

	err = paymentRepository.InsertPayment(payment)
	if err != nil {
		prepareFailureHeader(writer, request, err)
		return
	}

	writeHeaderLocation(writer, request, payment.ID)
	prepareSuccessHeader(writer, http.StatusCreated)
}

func updatePaymentEndpoint(writer http.ResponseWriter, request *http.Request) {
	var payment, err = decodeAndValidatePayment(request, false)
	if err != nil {
		prepareFailureHeader(writer, request, err)
		return
	}

	err = paymentRepository.UpdatePayment(payment)
	if err != nil {
		prepareFailureHeader(writer, request, err)
		return
	}

	writeHeaderLocation(writer, request, payment.ID)
	prepareSuccessHeader(writer, http.StatusOK)
}

func deletePaymentEndpoint(writer http.ResponseWriter, request *http.Request) {
	paymentID := mux.Vars(request)["id"]

	err := paymentRepository.DeletePayment(paymentID)
	if err != nil {
		prepareFailureHeader(writer, request, err)
		return
	}

	prepareSuccessHeader(writer, http.StatusOK)
}

func getPaymentEndpoint(writer http.ResponseWriter, request *http.Request) {
	paymentID := mux.Vars(request)["id"]

	payment, err := paymentRepository.GetPayment(paymentID)

	if err != nil {
		prepareFailureHeader(writer, request, err)
		return
	}

	prepareSuccessHeader(writer, http.StatusOK)

	links := Links{
		Self:   prepareFullPaymentURL(request.Host, getPaymentPath, paymentID),
		Update: prepareFullPaymentURL(request.Host, updatePaymentPath, ""),
		Delete: prepareFullPaymentURL(request.Host, deletePaymentPath, paymentID)}

	result := PaymentResult{payment, links}
	_ = json.NewEncoder(writer).Encode(result)
}

func getAllPaymentsEndpoint(writer http.ResponseWriter, request *http.Request) {
	payments, err := paymentRepository.GetAllPayments()

	if err != nil {
		prepareFailureHeader(writer, request, err)
		return
	}

	prepareSuccessHeader(writer, http.StatusOK)

	links := Links{
		Self: prepareFullPaymentURL(request.Host, getAllPaymentsPath, "")}

	result := PaymentListResult{payments, links}
	_ = json.NewEncoder(writer).Encode(result)
}
