package service

import "context"

// PaymentService 负责支付流程的编排。
type PaymentService interface {
	HandleWeChatCallback(ctx context.Context, payload []byte) error
}

type paymentService struct {
	deps         Dependencies
	orderService OrderService
}

// NewPaymentService 在订单服务基础上装配支付能力。
func NewPaymentService(deps Dependencies, orderService OrderService) PaymentService {
	return &paymentService{deps: deps, orderService: orderService}
}

func (s *paymentService) HandleWeChatCallback(ctx context.Context, payload []byte) error {
	s.deps.Logger.Printf("processing wechat callback size=%d", len(payload))

	result, err := s.deps.Payment.HandleCallback(ctx, payload)
	if err != nil {
		return err
	}

	if result.Success {
		if err := s.orderService.MarkPaid(ctx, result.OrderID); err != nil {
			return err
		}
	}

	return nil
}
