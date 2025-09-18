package service

import (
	"context"

	"convenienceStore/internal/model"
)

// ProductService exposes catalog operations.
type ProductService interface {
	ListProducts(ctx context.Context) ([]model.Product, error)
	GetProduct(ctx context.Context, productID string) (*model.Product, error)
	ValidateInventory(ctx context.Context, productID string, quantity int) (bool, error)
}

type productService struct {
	deps Dependencies
}

// NewProductService creates a new ProductService implementation.
func NewProductService(deps Dependencies) ProductService {
	return &productService{deps: deps}
}

func (s *productService) ListProducts(ctx context.Context) ([]model.Product, error) {
	s.deps.Logger.Println("listing products")
	return []model.Product{
		{
			ID:          "sku-demo",
			Name:        "Energy Drink",
			Description: "Refreshing energy beverage",
			Price:       4.5,
			Stock:       100,
			Tags:        []string{"drink", "energy"},
		},
	}, nil
}

func (s *productService) GetProduct(ctx context.Context, productID string) (*model.Product, error) {
	s.deps.Logger.Printf("fetching product %s", productID)
	return &model.Product{
		ID:          productID,
		Name:        "Sample Product",
		Description: "Product details placeholder",
		Price:       9.9,
		Stock:       50,
		Tags:        []string{"sample"},
	}, nil
}

func (s *productService) ValidateInventory(ctx context.Context, productID string, quantity int) (bool, error) {
	s.deps.Logger.Printf("validating stock product=%s qty=%d", productID, quantity)
	return quantity >= 0 && quantity <= 50, nil
}
