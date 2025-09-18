package service

import (
	"context"

	"convenienceStore/internal/model"
)

// DeliveryService 负责管理发货流程。
type DeliveryService interface {
	BindAddress(ctx context.Context, address *model.Address) error
	ShipOrder(ctx context.Context, orderID, carrier, trackNo string) error
}

type deliveryService struct {
	deps         Dependencies
	orderService OrderService
}

// NewDeliveryService 创建 DeliveryService 的实现。
func NewDeliveryService(deps Dependencies, orderService OrderService) DeliveryService {
	return &deliveryService{deps: deps, orderService: orderService}
}

func (s *deliveryService) BindAddress(ctx context.Context, address *model.Address) error {
	s.deps.Logger.Printf("binding delivery address %s for user %s", address.ID, address.UserID)
	return nil
}

func (s *deliveryService) ShipOrder(ctx context.Context, orderID, carrier, trackNo string) error {
	s.deps.Logger.Printf("shipping order %s with carrier=%s track=%s", orderID, carrier, trackNo)
	if err := s.orderService.ShipOrder(ctx, orderID); err != nil {
		return err
	}
	return nil
}
