package handler

import "convenienceStore/internal/service"

// Handlers 汇集各领域的 HTTP 处理器。
type Handlers struct {
	User     *UserHandler
	Product  *ProductHandler
	Cart     *CartHandler
	Order    *OrderHandler
	Payment  *PaymentHandler
	Delivery *DeliveryHandler
}

// NewHandlers 基于服务层依赖初始化所有处理器实例。
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
