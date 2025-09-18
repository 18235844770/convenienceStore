package service

import (
	"log"

	"convenienceStore/pkg/config"
	"convenienceStore/pkg/payment"
)

// Dependencies bundles the cross-cutting concerns required by the service layer.
type Dependencies struct {
	Config  *config.AppConfig
	Logger  *log.Logger
	Payment payment.WeChatClient
}

// Services exposes high level service singletons for consumers.
type Services struct {
	User     UserService
	Product  ProductService
	Cart     CartService
	Order    OrderService
	Payment  PaymentService
	Delivery DeliveryService
}

// NewServices wires the service layer graph.
func NewServices(deps Dependencies) Services {
	orderService := NewOrderService(deps)

	return Services{
		User:     NewUserService(deps),
		Product:  NewProductService(deps),
		Cart:     NewCartService(deps),
		Order:    orderService,
		Payment:  NewPaymentService(deps, orderService),
		Delivery: NewDeliveryService(deps, orderService),
	}
}
