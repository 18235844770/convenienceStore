package service

import (
	"context"
	"fmt"
	"time"

	"convenienceStore/internal/model"
	"convenienceStore/pkg/payment"
)

// OrderService coordinates order lifecycle operations.
type OrderService interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	GetOrder(ctx context.Context, orderID string) (*model.Order, error)
	PayOrder(ctx context.Context, orderID string) (*model.PaymentIntent, error)
	CancelOrder(ctx context.Context, orderID string) error
	ShipOrder(ctx context.Context, orderID string) error
	CompleteOrder(ctx context.Context, orderID string) error
	MarkPaid(ctx context.Context, orderID string) error
}

type orderService struct {
	deps Dependencies
}

// NewOrderService constructs the default OrderService implementation.
func NewOrderService(deps Dependencies) OrderService {
	return &orderService{deps: deps}
}

func (s *orderService) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	s.deps.Logger.Printf("creating order for user %s", order.UserID)
	order.ID = fmt.Sprintf("order-%d", time.Now().UnixNano())
	order.Status = model.OrderStatusPendingPayment
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	return order, nil
}

func (s *orderService) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	s.deps.Logger.Printf("fetching order %s", orderID)
	return &model.Order{
		ID:        orderID,
		UserID:    "user-demo",
		Status:    model.OrderStatusPaid,
		Total:     99.9,
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now(),
		Items: []model.OrderItem{
			{ProductID: "sku-demo", Quantity: 1, Price: 99.9},
		},
	}, nil
}

func (s *orderService) PayOrder(ctx context.Context, orderID string) (*model.PaymentIntent, error) {
	s.deps.Logger.Printf("initiating payment for order %s", orderID)
	resp, err := s.deps.Payment.CreateOrder(ctx, payment.OrderRequest{
		OrderID: orderID,
		Amount:  9990,
		Subject: "Convenience Store Order",
	})
	if err != nil {
		return nil, err
	}

	return &model.PaymentIntent{
		OrderID:     orderID,
		Provider:    "wechat",
		Credentials: resp.ClientConfig(),
	}, nil
}

func (s *orderService) CancelOrder(ctx context.Context, orderID string) error {
	s.deps.Logger.Printf("cancelling order %s", orderID)
	return nil
}

func (s *orderService) ShipOrder(ctx context.Context, orderID string) error {
	s.deps.Logger.Printf("shipping order %s", orderID)
	return nil
}

func (s *orderService) CompleteOrder(ctx context.Context, orderID string) error {
	s.deps.Logger.Printf("completing order %s", orderID)
	return nil
}

func (s *orderService) MarkPaid(ctx context.Context, orderID string) error {
	s.deps.Logger.Printf("marking order %s as paid", orderID)
	return nil
}
