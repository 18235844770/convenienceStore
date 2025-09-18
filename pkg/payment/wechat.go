package payment

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Config 定义接入微信支付所需的凭据。
type Config struct {
	AppID     string `mapstructure:"app_id"`
	MchID     string `mapstructure:"mch_id"`
	APIKey    string `mapstructure:"api_key"`
	NotifyURL string `mapstructure:"notify_url"`
}

// OrderRequest 描述创建支付订单所需的最小请求载荷。
type OrderRequest struct {
	OrderID string
	Amount  int64
	Subject string
}

// OrderResponse 封装支付服务方返回的、用于客户端的字段。
type OrderResponse struct {
	PrepayID  string
	NonceStr  string
	Timestamp string
	Signature string
}

// ClientConfig 将网关响应转换为前端易用的键值对。
func (r OrderResponse) ClientConfig() map[string]string {
	return map[string]string{
		"prepay_id": r.PrepayID,
		"nonce_str": r.NonceStr,
		"timestamp": r.Timestamp,
		"signature": r.Signature,
	}
}

// CallbackResult 汇总支付服务回调的核心字段。
type CallbackResult struct {
	OrderID string
	Success bool
}

// WeChatClient 抽象出对微信支付的调用接口。
type WeChatClient interface {
	CreateOrder(ctx context.Context, request OrderRequest) (OrderResponse, error)
	HandleCallback(ctx context.Context, payload []byte) (CallbackResult, error)
}

type weChatClient struct {
	cfg    Config
	logger *log.Logger
}

// NewWeChatClient 提供适合早期开发使用的桩实现。
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
