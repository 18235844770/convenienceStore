package payment

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Config defines credentials necessary to interact with WeChat Pay.
type Config struct {
	AppID     string `mapstructure:"app_id"`
	MchID     string `mapstructure:"mch_id"`
	APIKey    string `mapstructure:"api_key"`
	NotifyURL string `mapstructure:"notify_url"`
}

// OrderRequest describes the minimal payload required to create a payment order.
type OrderRequest struct {
	OrderID string
	Amount  int64
	Subject string
}

// OrderResponse encapsulates fields returned by the payment provider for client-side usage.
type OrderResponse struct {
	PrepayID  string
	NonceStr  string
	Timestamp string
	Signature string
}

// ClientConfig converts gateway response into a frontend friendly key-value map.
func (r OrderResponse) ClientConfig() map[string]string {
	return map[string]string{
		"prepay_id": r.PrepayID,
		"nonce_str": r.NonceStr,
		"timestamp": r.Timestamp,
		"signature": r.Signature,
	}
}

// CallbackResult summarizes the provider callback payload.
type CallbackResult struct {
	OrderID string
	Success bool
}

// WeChatClient abstracts access to WeChat Pay.
type WeChatClient interface {
	CreateOrder(ctx context.Context, request OrderRequest) (OrderResponse, error)
	HandleCallback(ctx context.Context, payload []byte) (CallbackResult, error)
}

type weChatClient struct {
	cfg    Config
	logger *log.Logger
}

// NewWeChatClient produces a stubbed implementation suitable for early development.
func NewWeChatClient(cfg Config, logger *log.Logger) WeChatClient {
	return &weChatClient{cfg: cfg, logger: logger}
}

func (c *weChatClient) CreateOrder(ctx context.Context, request OrderRequest) (OrderResponse, error) {
	c.logger.Printf("wechat create order order_id=%s amount=%d", request.OrderID, request.Amount)
	return OrderResponse{
		PrepayID:  fmt.Sprintf("wx_prepay_%s", request.OrderID),
		NonceStr:  "nonce-demo",
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Signature: "signature-demo",
	}, nil
}

func (c *weChatClient) HandleCallback(ctx context.Context, payload []byte) (CallbackResult, error) {
	c.logger.Printf("wechat callback received bytes=%d", len(payload))
	return CallbackResult{
		OrderID: "order-demo",
		Success: true,
	}, nil
}
