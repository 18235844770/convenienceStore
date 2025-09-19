package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"convenienceStore/internal/model"
	"convenienceStore/pkg/uid"
)

// CartService 负责处理购物车相关操作。
type CartService interface {
	ListItems(ctx context.Context, userID string) ([]model.CartItem, error)
	AddItem(ctx context.Context, item *model.CartItem) error
	UpdateItem(ctx context.Context, item *model.CartItem) error
	RemoveItem(ctx context.Context, itemID string) error
}

var errCartDBUnavailable = errors.New("cart service database is not configured")

type cartService struct {
	deps Dependencies
}

// NewCartService 返回购物车服务的基础实现。
func NewCartService(deps Dependencies) CartService {
	return &cartService{deps: deps}
}

func (s *cartService) ListItems(ctx context.Context, userID string) ([]model.CartItem, error) {
	if s.deps.DB == nil {
		return nil, errCartDBUnavailable
	}
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	const query = `SELECT id, user_id, product_id, quantity, selected, price FROM cart_items WHERE user_id = ? ORDER BY updated_at DESC`
	rows, err := s.deps.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.CartItem
	for rows.Next() {
		var item model.CartItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.Selected, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *cartService) AddItem(ctx context.Context, item *model.CartItem) error {
	if s.deps.DB == nil {
		return errCartDBUnavailable
	}
	if item == nil {
		return errors.New("cart item is nil")
	}
	if item.UserID == "" || item.ProductID == "" {
		return errors.New("user id and product id are required")
	}
	if item.ID == "" {
		item.ID = uid.New("cart_")
	}

	if item.Price == 0 {
		const priceQuery = `SELECT price FROM products WHERE id = ?`
		if err := s.deps.DB.QueryRowContext(ctx, priceQuery, item.ProductID).Scan(&item.Price); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("product %s not found", item.ProductID)
			}
			return err
		}
	}

	const stmt = `INSERT INTO cart_items (id, user_id, product_id, quantity, selected, price) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.deps.DB.ExecContext(ctx, stmt, item.ID, item.UserID, item.ProductID, item.Quantity, item.Selected, item.Price)
	return err
}

func (s *cartService) UpdateItem(ctx context.Context, item *model.CartItem) error {
	if s.deps.DB == nil {
		return errCartDBUnavailable
	}
	if item == nil {
		return errors.New("cart item is nil")
	}
	if item.ID == "" {
		return errors.New("cart item id is required")
	}

	const stmt = `UPDATE cart_items SET quantity = ?, selected = ?, price = ?, updated_at = NOW() WHERE id = ?`
	res, err := s.deps.DB.ExecContext(ctx, stmt, item.Quantity, item.Selected, item.Price, item.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("cart item %s not found", item.ID)
	}

	return nil
}

func (s *cartService) RemoveItem(ctx context.Context, itemID string) error {
	if s.deps.DB == nil {
		return errCartDBUnavailable
	}
	if itemID == "" {
		return errors.New("cart item id is required")
	}

	const stmt = `DELETE FROM cart_items WHERE id = ?`
	res, err := s.deps.DB.ExecContext(ctx, stmt, itemID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("cart item %s not found", itemID)
	}

	return nil
}
