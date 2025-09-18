package model

// User represents a shopper authenticated through WeChat.
type User struct {
	ID            string `json:"id"`
	WeChatOpenID  string `json:"wechat_open_id"`
	Nickname      string `json:"nickname"`
	AvatarURL     string `json:"avatar_url"`
	Phone         string `json:"phone"`
	DefaultAddrID string `json:"default_address_id"`
}

// Address captures a shipping destination for an order.
type Address struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Recipient  string `json:"recipient"`
	Phone      string `json:"phone"`
	Province   string `json:"province"`
	City       string `json:"city"`
	District   string `json:"district"`
	Detail     string `json:"detail"`
	PostalCode string `json:"postal_code"`
	IsDefault  bool   `json:"is_default"`
}
