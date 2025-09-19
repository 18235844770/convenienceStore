package service

import (
	"database/sql"
	"log"

	"convenienceStore/pkg/config"
	"convenienceStore/pkg/payment"
)

// Dependencies 汇集服务层所需的横切依赖。
type Dependencies struct {
	Config  *config.AppConfig
	Logger  *log.Logger
	DB      *sql.DB
	Payment payment.WeChatClient
}

// Services 对外暴露各领域的服务单例。
type Services struct {
	User     UserService
	Product  ProductService
	Cart     CartService
	Order    OrderService
	Payment  PaymentService
	Delivery DeliveryService
}

// NewServices 负责装配整个服务层依赖关系。
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
