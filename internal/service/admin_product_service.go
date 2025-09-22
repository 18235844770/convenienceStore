package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"convenienceStore/internal/model"
	"convenienceStore/pkg/uid"
)

// AdminProductPayload represents the editable attributes of a product.
type AdminProductPayload struct {
	Name        string
	Description string
	Price       float64
	Stock       int
	Tags        []string
	Images      []string
	IsActive    *bool
}

// AdminProductService exposes management operations for products.
type AdminProductService interface {
	ListProducts(ctx context.Context, status *bool) ([]model.Product, error)
	GetProduct(ctx context.Context, productID string) (*model.Product, error)
	CreateProduct(ctx context.Context, payload AdminProductPayload) (*model.Product, error)
	UpdateProduct(ctx context.Context, productID string, payload AdminProductPayload) (*model.Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	SetProductStatus(ctx context.Context, productID string, isActive bool) error
}

var errAdminProductDBUnavailable = errors.New("admin product service database is not configured")

type adminProductService struct {
	deps Dependencies
}

// NewAdminProductService creates an AdminProductService implementation.
func NewAdminProductService(deps Dependencies) AdminProductService {
	return &adminProductService{deps: deps}
}

func (s *adminProductService) ListProducts(ctx context.Context, status *bool) ([]model.Product, error) {
	if s.deps.DB == nil {
		return nil, errAdminProductDBUnavailable
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
		product, err := scanProductRow(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *adminProductService) GetProduct(ctx context.Context, productID string) (*model.Product, error) {
	if s.deps.DB == nil {
		return nil, errAdminProductDBUnavailable
	}

	const query = `SELECT id, name, description, price, stock, tags, images, is_active FROM products WHERE id = ?`
	row := s.deps.DB.QueryRowContext(ctx, query, productID)
	product, err := scanProductRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product %s not found", productID)
		}
		return nil, err
	}

	return product, nil
}

func (s *adminProductService) CreateProduct(ctx context.Context, payload AdminProductPayload) (*model.Product, error) {
	if s.deps.DB == nil {
		return nil, errAdminProductDBUnavailable
	}

	if err := validateAdminProductPayload(payload); err != nil {
		return nil, err
	}

	id := uid.New("prd_")
	tagsJSON, err := stringSliceToJSONArg(payload.Tags)
	if err != nil {
		return nil, err
	}
	imagesJSON, err := stringSliceToJSONArg(payload.Images)
	if err != nil {
		return nil, err
	}

	isActive := true
	if payload.IsActive != nil {
		isActive = *payload.IsActive
	}

	const query = `INSERT INTO products (id, name, description, price, stock, tags, images, is_active) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := s.deps.DB.ExecContext(ctx, query, id, payload.Name, payload.Description, payload.Price, payload.Stock, tagsJSON, imagesJSON, isActive); err != nil {
		return nil, err
	}

	return s.GetProduct(ctx, id)
}

func (s *adminProductService) UpdateProduct(ctx context.Context, productID string, payload AdminProductPayload) (*model.Product, error) {
	if s.deps.DB == nil {
		return nil, errAdminProductDBUnavailable
	}

	if err := validateAdminProductPayload(payload); err != nil {
		return nil, err
	}

	tagsJSON, err := stringSliceToJSONArg(payload.Tags)
	if err != nil {
		return nil, err
	}
	imagesJSON, err := stringSliceToJSONArg(payload.Images)
	if err != nil {
		return nil, err
	}

	query := `UPDATE products SET name = ?, description = ?, price = ?, stock = ?, tags = ?, images = ?`
	args := []any{payload.Name, payload.Description, payload.Price, payload.Stock, tagsJSON, imagesJSON}
	if payload.IsActive != nil {
		query += `, is_active = ?`
		args = append(args, *payload.IsActive)
	}
	query += ` WHERE id = ?`
	args = append(args, productID)

	result, err := s.deps.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, fmt.Errorf("product %s not found", productID)
	}

	return s.GetProduct(ctx, productID)
}

func (s *adminProductService) DeleteProduct(ctx context.Context, productID string) error {
	if s.deps.DB == nil {
		return errAdminProductDBUnavailable
	}

	const query = `DELETE FROM products WHERE id = ?`
	result, err := s.deps.DB.ExecContext(ctx, query, productID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("product %s not found", productID)
	}

	return nil
}

func (s *adminProductService) SetProductStatus(ctx context.Context, productID string, isActive bool) error {
	if s.deps.DB == nil {
		return errAdminProductDBUnavailable
	}

	const query = `UPDATE products SET is_active = ? WHERE id = ?`
	result, err := s.deps.DB.ExecContext(ctx, query, isActive, productID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("product %s not found", productID)
	}

	return nil
}

func validateAdminProductPayload(payload AdminProductPayload) error {
	if payload.Name == "" {
		return errors.New("product name is required")
	}
	if payload.Price < 0 {
		return errors.New("product price cannot be negative")
	}
	if payload.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}
	return nil
}

func scanProductRow(scanner interface {
	Scan(dest ...any) error
}) (*model.Product, error) {
	var (
		p        model.Product
		tags     sql.NullString
		images   sql.NullString
		isActive bool
	)

	if err := scanner.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &tags, &images, &isActive); err != nil {
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
