package model

import "time"

// OrderStatus captures the lifecycle of an order.
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusShipped        OrderStatus = "SHIPPED"
	OrderStatusCompleted      OrderStatus = "COMPLETED"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

var (
	_ = OrderStatusShipped
	_ = OrderStatusCompleted
	_ = OrderStatusCancelled
)

// OrderItem represents a product purchased within an order.
type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Order contains the purchase information and state transitions.
type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Status    OrderStatus `json:"status"`
	Total     float64     `json:"total"`
	AddressID string      `json:"address_id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
