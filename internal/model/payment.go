package model

// PaymentIntent contains client-side configuration values for payment initialization.
type PaymentIntent struct {
	OrderID     string            `json:"order_id"`
	Provider    string            `json:"provider"`
	Credentials map[string]string `json:"credentials"`
}
