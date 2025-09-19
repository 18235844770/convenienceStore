package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"convenienceStore/internal/model"
)

// ProductService 提供商品目录相关的业务能力。
type ProductService interface {
	ListProducts(ctx context.Context) ([]model.Product, error)
	GetProduct(ctx context.Context, productID string) (*model.Product, error)
	ValidateInventory(ctx context.Context, productID string, quantity int) (bool, error)
}

var errProductDBUnavailable = errors.New("product service database is not configured")

type productService struct {
	deps Dependencies
}

// NewProductService 创建新的 ProductService 实现。
func NewProductService(deps Dependencies) ProductService {
	return &productService{deps: deps}
}

func (s *productService) ListProducts(ctx context.Context) ([]model.Product, error) {
	if s.deps.DB == nil {
		return nil, errProductDBUnavailable
	}

	const query = `SELECT id, name, description, price, stock, tags FROM products ORDER BY updated_at DESC`
	rows, err := s.deps.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		var tags sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &tags); err != nil {
			return nil, err
		}

		if tags.Valid && tags.String != "" {
			var parsed []string
			if err := json.Unmarshal([]byte(tags.String), &parsed); err == nil {
				p.Tags = parsed
			}
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *productService) GetProduct(ctx context.Context, productID string) (*model.Product, error) {
	if s.deps.DB == nil {
		return nil, errProductDBUnavailable
	}

	const query = `SELECT id, name, description, price, stock, tags FROM products WHERE id = ?`
	var p model.Product
	var tags sql.NullString
	if err := s.deps.DB.QueryRowContext(ctx, query, productID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &tags); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product %s not found", productID)
		}
		return nil, err
	}

	if tags.Valid && tags.String != "" {
		var parsed []string
		if err := json.Unmarshal([]byte(tags.String), &parsed); err == nil {
			p.Tags = parsed
		}
	}

	return &p, nil
}

func (s *productService) ValidateInventory(ctx context.Context, productID string, quantity int) (bool, error) {
	if s.deps.DB == nil {
		return false, errProductDBUnavailable
	}

	if quantity <= 0 {
		return false, nil
	}

	const query = `SELECT stock FROM products WHERE id = ?`
	var stock int
	if err := s.deps.DB.QueryRowContext(ctx, query, productID).Scan(&stock); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("product %s not found", productID)
		}
		return false, err
	}

	return quantity <= stock, nil
}
