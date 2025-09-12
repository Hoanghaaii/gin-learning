package entities

import "time"

type OrderPaymentStatus string

const (
	OrderPaymentPending OrderPaymentStatus = "order_payment_pending"
	OrderPaymentPaid    OrderPaymentStatus = "order_payment_paid"
	OrderPaymentFailed  OrderPaymentStatus = "order_payment_failed"
)

type OrderStatus string

const (
	OrderStatusPending OrderStatus = "order_status_pending"
	OrderStatusFinish  OrderStatus = "order_status_finish"
	OrderStatusFailed  OrderStatus = "order_status_failed"
)

type Order struct {
	ID            uint
	UserID        uint
	EventID       uint
	FinalPrice    float64
	Quantity      uint
	TotalAmount   float64
	PaymentStatus OrderPaymentStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
