package model

// CartItem 表示用户购物车中的单个商品项。
type CartItem struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Selected  bool    `json:"selected"`
	Price     float64 `json:"price"`
}
