package entities

import "time"

type PaymentMethod string

const (
	MethodVNPAY PaymentMethod = "vnpay"
	MethodMOMO  PaymentMethod = "momo"
	MethodCOD   PaymentMethod = "cod"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "payment_pending"
	PaymentStatusSuccess  PaymentStatus = "payment_success"
	PaymentStatusFailed   PaymentStatus = "payment_failed"
	PaymentStatusRefunded PaymentStatus = "payment_refunded"
)

type Payment struct {
	ID            uint
	OrderID       uint
	Method        PaymentMethod
	Status        PaymentStatus
	TransactionID string
	CreatedAt     time.Time
}
