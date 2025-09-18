package service

import "context"

// PaymentService handles payment orchestration tasks.
type PaymentService interface {
	HandleWeChatCallback(ctx context.Context, payload []byte) error
}

type paymentService struct {
	deps         Dependencies
	orderService OrderService
}

// NewPaymentService wires payment functionality on top of the order service.
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
