package service

import (
	"context"

	"convenienceStore/internal/model"
)

// CartService manages shopping cart operations.
type CartService interface {
	ListItems(ctx context.Context, userID string) ([]model.CartItem, error)
	AddItem(ctx context.Context, item *model.CartItem) error
	UpdateItem(ctx context.Context, item *model.CartItem) error
	RemoveItem(ctx context.Context, itemID string) error
}

type cartService struct {
	deps Dependencies
}

// NewCartService returns a baseline cart service.
func NewCartService(deps Dependencies) CartService {
	return &cartService{deps: deps}
}

func (s *cartService) ListItems(ctx context.Context, userID string) ([]model.CartItem, error) {
	s.deps.Logger.Printf("listing cart items for user %s", userID)
	return []model.CartItem{
		{
			ID:        "cart-item-demo",
			UserID:    userID,
			ProductID: "sku-demo",
			Quantity:  2,
			Selected:  true,
		},
	}, nil
}

func (s *cartService) AddItem(ctx context.Context, item *model.CartItem) error {
	s.deps.Logger.Printf("adding cart item %s for user %s", item.ProductID, item.UserID)
	return nil
}

func (s *cartService) UpdateItem(ctx context.Context, item *model.CartItem) error {
	s.deps.Logger.Printf("updating cart item %s", item.ID)
	return nil
}

func (s *cartService) RemoveItem(ctx context.Context, itemID string) error {
	s.deps.Logger.Printf("removing cart item %s", itemID)
	return nil
}
