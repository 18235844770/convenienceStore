package model

// PaymentIntent 包含客户端初始化支付所需的配置。
type PaymentIntent struct {
	OrderID     string            `json:"order_id"`
	Provider    string            `json:"provider"`
	Credentials map[string]string `json:"credentials"`
}
