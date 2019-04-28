package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

var paymentRepository PaymentRepository

func SetPaymentRepository(repository PaymentRepository) {
	paymentRepository = repository
}

func PrepareSuccessHeader(writer http.ResponseWriter, statusCode int) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(statusCode)
}

func PrepareFailureHeader(writer http.ResponseWriter, err error) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")

	switch err.(type) {
	case *PersistenceError:
		writer.WriteHeader(http.StatusInternalServerError)
		return
	case *NotFoundError:
		writer.WriteHeader(http.StatusNotFound)
		return
	case *VersionConflictError:
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

func DecodeAndValidatePayment(request *http.Request, create bool) (payment Payment, err error) {
	err = json.NewDecoder(request.Body).Decode(&payment)
	if err != nil {
		return payment, err
	}

	if len(payment.OrganisationId) == 0 || (!create && len(payment.ID) == 0) {
		return payment, &InvalidPaymentError{payment}
	}

	return payment, nil
}

// TODO log errors
func CreatePaymentHandler(writer http.ResponseWriter, request *http.Request) {
	var payment, err = DecodeAndValidatePayment(request, true)
	if err != nil {
		PrepareFailureHeader(writer, err)
		return
	}

	newUuid, _ := uuid.NewUUID()
	payment.ID = newUuid.String()
	payment.Version = 0

	err = paymentRepository.InsertPayment(payment)
	if err != nil {
		PrepareFailureHeader(writer, err)
		return
	}

	PrepareSuccessHeader(writer, http.StatusCreated)
	writer.Header().Set("Location", string(request.URL.Path)+"/"+payment.ID)
}

// TODO log errors
func UpdatePaymentHandler(writer http.ResponseWriter, request *http.Request) {
	var payment, err = DecodeAndValidatePayment(request, false)
	if err != nil {
		PrepareFailureHeader(writer, err)
		return
	}

	err = paymentRepository.UpdatePayment(payment)
	if err != nil {
		PrepareFailureHeader(writer, err)
		return
	}

	PrepareSuccessHeader(writer, http.StatusOK)
	writer.Header().Set("Location", string(request.URL.Path)+"/"+payment.ID)
}

// TODO log errors
func DeletePaymentHandler(writer http.ResponseWriter, request *http.Request) {
	paymentId := mux.Vars(request)["id"]

	err := paymentRepository.DeletePayment(paymentId)
	if err != nil {
		PrepareFailureHeader(writer, err)
		return
	}

	PrepareSuccessHeader(writer, http.StatusOK)
}

// TODO log errors
func GetPaymentHandler(writer http.ResponseWriter, request *http.Request) {
	paymentId := mux.Vars(request)["id"]

	payment, err := paymentRepository.GetPayment(paymentId)

	if err != nil {
		PrepareFailureHeader(writer, err)
		return
	}

	PrepareSuccessHeader(writer, http.StatusOK)
	json.NewEncoder(writer).Encode(payment)
}

// TODO log errors
func GetAllPaymentsHandler(writer http.ResponseWriter, request *http.Request) {
	payments, err := paymentRepository.GetAllPayments()

	if err != nil {
		PrepareFailureHeader(writer, err)
		return
	}

	PrepareSuccessHeader(writer, http.StatusOK)
	json.NewEncoder(writer).Encode(payments)
}