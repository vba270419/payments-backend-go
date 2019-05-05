package main

import "fmt"

// A PersistenceError is an error type when any query or update to a storage failed due to a technical reason
type PersistenceError struct {
}

func (e PersistenceError) Error() string {
	return "Database query failed"
}

// A PaymentNotFoundError is an error type when Payment for a given paymentID can not be found in the storage
type PaymentNotFoundError struct {
	paymentID string
}

func (e PaymentNotFoundError) Error() string {
	return fmt.Sprintf("Payment '%s' not found", e.paymentID)
}

// A PaymentVersionConflictError is an error type when Payment object can not be persisted to a storage due to a version conflict
type PaymentVersionConflictError struct {
	paymentID string
	version   int
}

func (e PaymentVersionConflictError) Error() string {
	return fmt.Sprintf("Payment '%s' with version '%d' can not be updated", e.paymentID, e.version)
}

// An InvalidPaymentError is an error type when given Payment object has invalid or inconsistent data
type InvalidPaymentError struct {
	payment Payment
}

func (e InvalidPaymentError) Error() string {
	return fmt.Sprintf("Payment has invalid format %+v\n", e.payment)
}
