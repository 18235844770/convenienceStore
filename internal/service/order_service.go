package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"convenienceStore/internal/model"
	"convenienceStore/pkg/payment"
	"convenienceStore/pkg/uid"
)

// OrderService 调度订单全生命周期的业务操作。
type OrderService interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	GetOrder(ctx context.Context, orderID string) (*model.Order, error)
	PayOrder(ctx context.Context, orderID string) (*model.PaymentIntent, error)
	CancelOrder(ctx context.Context, orderID string) error
	ShipOrder(ctx context.Context, orderID string) error
	CompleteOrder(ctx context.Context, orderID string) error
	MarkPaid(ctx context.Context, orderID string) error
}

var errOrderDBUnavailable = errors.New("order service database is not configured")

type orderService struct {
	deps Dependencies
}

// NewOrderService 构建默认的 OrderService 实现。
func NewOrderService(deps Dependencies) OrderService {
	return &orderService{deps: deps}
}

func (s *orderService) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	if s.deps.DB == nil {
		return nil, errOrderDBUnavailable
	}
	if order == nil {
		return nil, errors.New("order payload is nil")
	}
	if order.UserID == "" {
		return nil, errors.New("user id is required")
	}

	if order.ID == "" {
		order.ID = uid.New("ord_")
	}

	var total float64
	for i := range order.Items {
		item := &order.Items[i]
		if item.ProductID == "" {
			return nil, errors.New("order item product id is required")
		}
		if item.Quantity <= 0 {
			return nil, errors.New("order item quantity must be positive")
		}
		if item.Price == 0 {
			const priceQuery = `SELECT price FROM products WHERE id = ?`
			if err := s.deps.DB.QueryRowContext(ctx, priceQuery, item.ProductID).Scan(&item.Price); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, fmt.Errorf("product %s not found", item.ProductID)
				}
				return nil, err
			}
		}
		total += item.Price * float64(item.Quantity)
	}
	order.Total = total

	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	order.Status = model.OrderStatusPendingPayment

	tx, err := s.deps.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	const orderInsert = `INSERT INTO orders (id, user_id, status, total, address_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	if _, err = tx.ExecContext(ctx, orderInsert, order.ID, order.UserID, order.Status, order.Total, order.AddressID, order.CreatedAt, order.UpdatedAt); err != nil {
		return nil, err
	}

	const itemInsert = `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`
	for _, item := range order.Items {
		if _, err = tx.ExecContext(ctx, itemInsert, order.ID, item.ProductID, item.Quantity, item.Price); err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	if s.deps.DB == nil {
		return nil, errOrderDBUnavailable
	}
	if orderID == "" {
		return nil, errors.New("order id is required")
	}

	const orderQuery = `SELECT id, user_id, status, total, address_id, created_at, updated_at FROM orders WHERE id = ?`
	var order model.Order
	if err := s.deps.DB.QueryRowContext(ctx, orderQuery, orderID).Scan(&order.ID, &order.UserID, &order.Status, &order.Total, &order.AddressID, &order.CreatedAt, &order.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order %s not found", orderID)
		}
		return nil, err
	}

	const itemsQuery = `SELECT product_id, quantity, price FROM order_items WHERE order_id = ?`
	rows, err := s.deps.DB.QueryContext(ctx, itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *orderService) PayOrder(ctx context.Context, orderID string) (*model.PaymentIntent, error) {
	if s.deps.DB == nil {
		return nil, errOrderDBUnavailable
	}

	const amountQuery = `SELECT total FROM orders WHERE id = ?`
	var total float64
	if err := s.deps.DB.QueryRowContext(ctx, amountQuery, orderID).Scan(&total); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order %s not found", orderID)
		}
		return nil, err
	}

	amount := int64(math.Round(total * 100))
	resp, err := s.deps.Payment.CreateOrder(ctx, payment.OrderRequest{
		OrderID: orderID,
		Amount:  amount,
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
	return s.updateStatus(ctx, orderID, model.OrderStatusCancelled)
}

func (s *orderService) ShipOrder(ctx context.Context, orderID string) error {
	return s.updateStatus(ctx, orderID, model.OrderStatusShipped)
}

func (s *orderService) CompleteOrder(ctx context.Context, orderID string) error {
	return s.updateStatus(ctx, orderID, model.OrderStatusCompleted)
}

func (s *orderService) MarkPaid(ctx context.Context, orderID string) error {
	return s.updateStatus(ctx, orderID, model.OrderStatusPaid)
}

func (s *orderService) updateStatus(ctx context.Context, orderID string, status model.OrderStatus) error {
	if s.deps.DB == nil {
		return errOrderDBUnavailable
	}
	if orderID == "" {
		return errors.New("order id is required")
	}

	const stmt = `UPDATE orders SET status = ?, updated_at = NOW() WHERE id = ?`
	res, err := s.deps.DB.ExecContext(ctx, stmt, status, orderID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("order %s not found", orderID)
	}

	return nil
}
