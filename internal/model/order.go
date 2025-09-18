package model

import "time"

// OrderStatus 描述订单生命周期。
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

// OrderItem 表示订单中购买的单件商品。
type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Order 包含订单的核心信息及状态流转。
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
