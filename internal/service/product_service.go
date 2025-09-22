package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"convenienceStore/internal/model"
)

// ProductService exposes product catalog operations.
type ProductService interface {
	ListProducts(ctx context.Context, status *bool) ([]model.Product, error)
	GetProduct(ctx context.Context, productID string, status *bool) (*model.Product, error)
	ValidateInventory(ctx context.Context, productID string, quantity int) (bool, error)
}

var errProductDBUnavailable = errors.New("product service database is not configured")

type productService struct {
	deps Dependencies
}

// NewProductService creates a new ProductService implementation.
func NewProductService(deps Dependencies) ProductService {
	return &productService{deps: deps}
}

func parseStringArray(raw sql.NullString) []string {
	if !raw.Valid || raw.String == "" {
		return nil
	}

	var parsed []string
	if err := json.Unmarshal([]byte(raw.String), &parsed); err != nil {
		return nil
	}

	return parsed
}

func (s *productService) ListProducts(ctx context.Context, status *bool) ([]model.Product, error) {
	if s.deps.DB == nil {
		return nil, errProductDBUnavailable
	}

	query := `SELECT id, name, description, price, stock, tags, images, is_active FROM products`
	var args []any
	if status != nil {
		query += ` WHERE is_active = ?`
		args = append(args, *status)
	}
	query += ` ORDER BY updated_at DESC`

	rows, err := s.deps.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		var tags sql.NullString
		var images sql.NullString
		var isActive bool
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &tags, &images, &isActive); err != nil {
			return nil, err
		}

		if parsedTags := parseStringArray(tags); parsedTags != nil {
			p.Tags = parsedTags
		}

		if parsedImages := parseStringArray(images); parsedImages != nil {
			p.Images = parsedImages
		}

		p.IsActive = isActive

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *productService) GetProduct(ctx context.Context, productID string, status *bool) (*model.Product, error) {
	if s.deps.DB == nil {
		return nil, errProductDBUnavailable
	}

	query := `SELECT id, name, description, price, stock, tags, images, is_active FROM products WHERE id = ?`
	args := []any{productID}
	if status != nil {
		query += ` AND is_active = ?`
		args = append(args, *status)
	}

	var p model.Product
	var tags sql.NullString
	var images sql.NullString
	var isActive bool
	if err := s.deps.DB.QueryRowContext(ctx, query, args...).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &tags, &images, &isActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product %s not found", productID)
		}
		return nil, err
	}

	if parsedTags := parseStringArray(tags); parsedTags != nil {
		p.Tags = parsedTags
	}

	if parsedImages := parseStringArray(images); parsedImages != nil {
		p.Images = parsedImages
	}

	p.IsActive = isActive

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
