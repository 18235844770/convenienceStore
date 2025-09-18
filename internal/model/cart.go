package model

// CartItem represents a single product entry within a user cart.
type CartItem struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Selected  bool    `json:"selected"`
	Price     float64 `json:"price"`
}
