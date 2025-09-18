package handler

import "convenienceStore/internal/service"

// Handlers combines each domain specific HTTP handler.
type Handlers struct {
	User     *UserHandler
	Product  *ProductHandler
	Cart     *CartHandler
	Order    *OrderHandler
	Payment  *PaymentHandler
	Delivery *DeliveryHandler
}

// NewHandlers bootstraps handler instances from the service layer dependencies.
func NewHandlers(services service.Services) Handlers {
	return Handlers{
		User:     NewUserHandler(services.User),
		Product:  NewProductHandler(services.Product),
		Cart:     NewCartHandler(services.Cart),
		Order:    NewOrderHandler(services.Order),
		Payment:  NewPaymentHandler(services.Payment),
		Delivery: NewDeliveryHandler(services.Delivery),
	}
}
