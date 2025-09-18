package model

// User 表示通过微信完成认证的顾客。
type User struct {
	ID            string `json:"id"`
	WeChatOpenID  string `json:"wechat_open_id"`
	Nickname      string `json:"nickname"`
	AvatarURL     string `json:"avatar_url"`
	Phone         string `json:"phone"`
	DefaultAddrID string `json:"default_address_id"`
}

// Address 表示订单的收货目的地。
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
