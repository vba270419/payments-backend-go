package main

import "fmt"

type PersistenceError struct {
}

func (e PersistenceError) Error() string {
	return "Database query failed"
}

type NotFoundError struct {
	paymentId string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("Payment '%s' not found", e.paymentId)
}

type VersionConflictError struct {
	paymentId string
	version   int
}

func (e VersionConflictError) Error() string {
	return fmt.Sprintf("Payment '%s' with version '%d' can not be updated", e.paymentId, e.version)
}

type InvalidPaymentError struct {
	payment Payment
}

func (e InvalidPaymentError) Error() string {
	return fmt.Sprintf("Payment has invalid format %+v\n", e.payment)
}
