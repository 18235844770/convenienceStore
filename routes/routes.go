package routes

import (
	"convenienceStore/internal/handler"

	"github.com/gin-gonic/gin"
)

// HandlerSet 汇总应用所需的全量 HTTP 处理器。
type HandlerSet struct {
	User         *handler.UserHandler
	Product      *handler.ProductHandler
	AdminProduct *handler.AdminProductHandler
	Upload       *handler.UploadHandler
	Cart         *handler.CartHandler
	Order        *handler.OrderHandler
	Payment      *handler.PaymentHandler
	Delivery     *handler.DeliveryHandler
}

// RegisterRoutes 将各领域的路由绑定到对应处理器。
func RegisterRoutes(engine *gin.Engine, handlers HandlerSet) {
	api := engine.Group("/api")

	userGroup := api.Group("/users")
	userGroup.POST("/wechat/login", handlers.User.WeChatLogin)
	userGroup.POST("/bind", handlers.User.BindUser)
	userGroup.GET("/addresses", handlers.User.ListAddresses)
	userGroup.POST("/addresses", handlers.User.CreateAddress)
	userGroup.PUT("/addresses/:id", handlers.User.UpdateAddress)
	userGroup.DELETE("/addresses/:id", handlers.User.DeleteAddress)

	productGroup := api.Group("/products")
	productGroup.GET("", handlers.Product.ListProducts)
	productGroup.GET(":id", handlers.Product.GetProduct)
	productGroup.POST(":id/validate", handlers.Product.ValidateInventory)

	adminGroup := api.Group("/admin")
	adminProducts := adminGroup.Group("/products")
	adminProducts.GET("", handlers.AdminProduct.ListProducts)
	adminProducts.GET("/:id", handlers.AdminProduct.GetProduct)
	adminProducts.POST("", handlers.AdminProduct.CreateProduct)
	adminProducts.PUT("/:id", handlers.AdminProduct.UpdateProduct)
	adminProducts.DELETE("/:id", handlers.AdminProduct.DeleteProduct)
	adminProducts.PATCH("/:id/status", handlers.AdminProduct.SetProductStatus)

	adminGroup.POST("/uploads", handlers.Upload.UploadFile)

	cartGroup := api.Group("/cart")
	cartGroup.GET("", handlers.Cart.ListItems)
	cartGroup.POST("", handlers.Cart.AddItem)
	cartGroup.PUT(":id", handlers.Cart.UpdateItem)
	cartGroup.DELETE(":id", handlers.Cart.RemoveItem)

	orderGroup := api.Group("/orders")
	orderGroup.POST("", handlers.Order.CreateOrder)
	orderGroup.GET(":id", handlers.Order.GetOrder)
	orderGroup.POST(":id/pay", handlers.Order.PayOrder)
	orderGroup.POST(":id/cancel", handlers.Order.CancelOrder)
	orderGroup.POST(":id/ship", handlers.Order.ShipOrder)
	orderGroup.POST(":id/complete", handlers.Order.CompleteOrder)

	paymentGroup := api.Group("/payments")
	paymentGroup.POST("/wechat/callback", handlers.Payment.HandleWeChatCallback)

	deliveryGroup := api.Group("/delivery")
	deliveryGroup.POST("/bind-address", handlers.Delivery.BindAddress)
	deliveryGroup.POST("/ship-order", handlers.Delivery.ShipOrder)
}
