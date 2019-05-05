package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	. "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	successful      string = "successful"
	dbFailure       string = "dbFailure"
	notFound        string = "notFound"
	versionConflict string = "versionConflict"
)

type PaymentRepositoryMock struct {
	mock.Mock
	mode string
}

func (m *PaymentRepositoryMock) InsertPayment(payment Payment) (err error) {
	switch m.mode {
	case dbFailure:
		return &PersistenceError{}
	default:
		return
	}
}

func (m *PaymentRepositoryMock) UpdatePayment(payment Payment) (err error) {
	switch m.mode {
	case dbFailure:
		return &PersistenceError{}
	case notFound:
		return &NotFoundError{payment.ID}
	case versionConflict:
		return &VersionConflictError{payment.ID, payment.Version}
	default:
		return
	}
}

func (m *PaymentRepositoryMock) DeletePayment(paymentId string) (err error) {
	switch m.mode {
	case notFound:
		return &NotFoundError{paymentId}
	case dbFailure:
		return &PersistenceError{}
	default:
		return
	}
}

func (m *PaymentRepositoryMock) GetPayment(paymentId string) (payment Payment, err error) {
	payment = Payment{ID: paymentId, OrganisationID: "123", Version: 1}

	switch m.mode {
	case notFound:
		return payment, &NotFoundError{paymentId}
	case dbFailure:
		return payment, &PersistenceError{}
	default:
		return payment, nil
	}
}

func (m *PaymentRepositoryMock) GetAllPayments() (payments []Payment, err error) {
	payments = append(payments, Payment{ID: "1", OrganisationID: "123", Version: 1})
	payments = append(payments, Payment{ID: "2", OrganisationID: "456", Version: 2})
	payments = append(payments, Payment{ID: "3", OrganisationID: "789", Version: 3})

	switch m.mode {
	case dbFailure:
		return payments, &PersistenceError{}
	default:
		return payments, nil
	}
}

// Test payment creation handler

func TestCreatePaymentSuccessful(t *testing.T) {
	response := ServeHTTP(methodPost, createPaymentPath, MockPayment("", "123"), successful)

	Equal(t, 201, response.Code)
	Contains(t, response.Header().Get("Location"), "/v1/payments/get")
}

func TestCreatePaymentEmptyBody(t *testing.T) {
	response := ServeHTTP(methodPost, createPaymentPath, http.NoBody, successful)

	Equal(t, 400, response.Code)
}

func TestCreatePaymentInvalidFormat(t *testing.T) {
	response := ServeHTTP(methodPost, createPaymentPath, MockPayment("", ""), successful)

	Equal(t, 400, response.Code)
}

func TestCreatePaymentServerFailed(t *testing.T) {
	response := ServeHTTP(methodPost, createPaymentPath, MockPayment("", "123"), dbFailure)

	Equal(t, 500, response.Code)
}

// Test payment update handler

func TestUpdatePaymentSuccessful(t *testing.T) {
	response := ServeHTTP(methodPut, updatePaymentPath, MockPayment("1", "123"), successful)

	Equal(t, 200, response.Code)
	Contains(t, response.Header().Get("Location"), "/v1/payments/get/1")
}

func TestUpdatePaymentNotFound(t *testing.T) {
	response := ServeHTTP(methodPut, updatePaymentPath, MockPayment("1", "123"), notFound)

	Equal(t, 404, response.Code)
}

func TestUpdatePaymentVersionConflict(t *testing.T) {
	response := ServeHTTP(methodPut, updatePaymentPath, MockPayment("1", "123"), versionConflict)

	Equal(t, 409, response.Code)
}

func TestUpdatePaymentEmptyBody(t *testing.T) {
	response := ServeHTTP(methodPut, updatePaymentPath, http.NoBody, successful)

	Equal(t, 400, response.Code)
}

func TestUpdatePaymentNoId(t *testing.T) {
	response := ServeHTTP(methodPut, updatePaymentPath, MockPayment("", "123"), successful)

	Equal(t, 400, response.Code)
}

func TestUpdatePaymentServerFailed(t *testing.T) {
	response := ServeHTTP(methodPut, updatePaymentPath, MockPayment("1", "123"), dbFailure)

	Equal(t, 500, response.Code)
}

// Test payment delete handler

func TestDeletePaymentSuccessful(t *testing.T) {
	response := ServeHTTP(methodDelete, PreparePaymentURL(deletePaymentPath, "1"), http.NoBody, successful)

	Equal(t, 200, response.Code)
}

func TestDeletePaymentNotFound(t *testing.T) {
	response := ServeHTTP(methodDelete, PreparePaymentURL(deletePaymentPath, "1"), http.NoBody, notFound)

	Equal(t, 404, response.Code)
}

func TestDeletePaymentServerFailed(t *testing.T) {
	response := ServeHTTP(methodDelete, PreparePaymentURL(deletePaymentPath, "1"), http.NoBody, dbFailure)

	Equal(t, 500, response.Code)
}

// Test get payment handler
func TestGetPaymentSuccessful(t *testing.T) {
	response := ServeHTTP(methodGet, PreparePaymentURL(getPaymentPath, "2"), http.NoBody, successful)

	var paymentResult PaymentResult
	_ = json.NewDecoder(response.Body).Decode(&paymentResult)

	Equal(t, 200, response.Code)
	Equal(t, "2", paymentResult.Data.ID)
	Equal(t, "123", paymentResult.Data.OrganisationID)
	Equal(t, 1, paymentResult.Data.Version)
}

func TestGetPaymentNotFound(t *testing.T) {
	response := ServeHTTP(methodGet, PreparePaymentURL(getPaymentPath, "1"), http.NoBody, notFound)

	Equal(t, 404, response.Code)
}

func TestGetPaymentServerFailed(t *testing.T) {
	response := ServeHTTP(methodGet, PreparePaymentURL(getPaymentPath, "1"), http.NoBody, dbFailure)

	Equal(t, 500, response.Code)
}

// Test get all payments handler

func TestGetAllPaymentsSuccessful(t *testing.T) {
	response := ServeHTTP(methodGet, getAllPaymentsPath, http.NoBody, successful)

	var result PaymentListResult
	_ = json.NewDecoder(response.Body).Decode(&result)

	Equal(t, 200, response.Code)
	Equal(t, 3, len(result.Data))

	Equal(t, "1", result.Data[0].ID)
	Equal(t, "2", result.Data[1].ID)
	Equal(t, "3", result.Data[2].ID)

	Equal(t, "123", result.Data[0].OrganisationID)
	Equal(t, "456", result.Data[1].OrganisationID)
	Equal(t, "789", result.Data[2].OrganisationID)
}

func TestGetAllPaymentsServerFailed(t *testing.T) {
	response := ServeHTTP(methodGet, getAllPaymentsPath, http.NoBody, dbFailure)

	Equal(t, 500, response.Code)
}

// ---------------------------------------------------- //

func MockRouter(mode string) (router *mux.Router) {
	repository := new(PaymentRepositoryMock)
	repository.mode = mode

	SetPaymentRepository(repository)
	router = ConfigureRouter()
	return router
}

func MockPayment(id string, organisationId string) *bytes.Buffer {
	payment := &Payment{ID: id, OrganisationID: organisationId}
	jsonPayment, _ := json.Marshal(payment)
	return bytes.NewBuffer(jsonPayment)
}

func ServeHTTP(method string, url string, body io.Reader, mode string) *httptest.ResponseRecorder {
	router := MockRouter(mode)
	request, _ := http.NewRequest(method, url, body)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	return response
}
