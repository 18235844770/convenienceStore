package service

import (
	"context"

	"convenienceStore/internal/model"
)

// UserService defines behaviours around user accounts and addresses.
type UserService interface {
	WeChatLogin(ctx context.Context, code string) (*model.User, error)
	BindUser(ctx context.Context, user *model.User) error
	ListAddresses(ctx context.Context, userID string) ([]model.Address, error)
	CreateAddress(ctx context.Context, address *model.Address) error
	UpdateAddress(ctx context.Context, address *model.Address) error
	DeleteAddress(ctx context.Context, addressID string) error
}

type userService struct {
	deps Dependencies
}

// NewUserService provides a baseline implementation of UserService.
func NewUserService(deps Dependencies) UserService {
	return &userService{deps: deps}
}

func (s *userService) WeChatLogin(ctx context.Context, code string) (*model.User, error) {
	s.deps.Logger.Printf("wechat login attempt code=%s", code)
	return &model.User{
		ID:            "user-demo",
		WeChatOpenID:  "openid-demo",
		Nickname:      "New Shopper",
		AvatarURL:     "https://example.com/avatar.png",
		DefaultAddrID: "addr-demo",
	}, nil
}

func (s *userService) BindUser(ctx context.Context, user *model.User) error {
	s.deps.Logger.Printf("binding user %s", user.ID)
	return nil
}

func (s *userService) ListAddresses(ctx context.Context, userID string) ([]model.Address, error) {
	s.deps.Logger.Printf("listing addresses for user %s", userID)
	return []model.Address{
		{
			ID:         "addr-demo",
			UserID:     userID,
			Recipient:  "Demo User",
			Phone:      "18800000000",
			Province:   "Guangdong",
			City:       "Shenzhen",
			District:   "Nanshan",
			Detail:     "Tencent Building",
			IsDefault:  true,
			PostalCode: "518000",
		},
	}, nil
}

func (s *userService) CreateAddress(ctx context.Context, address *model.Address) error {
	s.deps.Logger.Printf("creating address for user %s", address.UserID)
	return nil
}

func (s *userService) UpdateAddress(ctx context.Context, address *model.Address) error {
	s.deps.Logger.Printf("updating address %s", address.ID)
	return nil
}

func (s *userService) DeleteAddress(ctx context.Context, addressID string) error {
	s.deps.Logger.Printf("deleting address %s", addressID)
	return nil
}
