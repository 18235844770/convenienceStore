package model

// Product represents an item that can be purchased.
type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	Tags        []string `json:"tags"`
	Images      []string `json:"images"`
	IsActive    bool     `json:"is_active"`
}
