package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"convenienceStore/internal/model"
	"convenienceStore/pkg/uid"
)

// UserService 定义与用户账号及地址相关的业务行为。
type UserService interface {
	WeChatLogin(ctx context.Context, code string) (*model.User, error)
	BindUser(ctx context.Context, user *model.User) error
	ListAddresses(ctx context.Context, userID string) ([]model.Address, error)
	CreateAddress(ctx context.Context, address *model.Address) error
	UpdateAddress(ctx context.Context, address *model.Address) error
	DeleteAddress(ctx context.Context, addressID string) error
}

var errUserDBUnavailable = errors.New("user service database is not configured")

type userService struct {
	deps Dependencies
}

// NewUserService 提供 UserService 的基础实现。
func NewUserService(deps Dependencies) UserService {
	return &userService{deps: deps}
}

func (s *userService) WeChatLogin(ctx context.Context, code string) (*model.User, error) {
	if s.deps.DB == nil {
		return nil, errUserDBUnavailable
	}
	if code == "" {
		return nil, errors.New("wechat auth code is required")
	}

	openID := code

	const query = `SELECT id, wechat_open_id, nickname, avatar_url, phone, default_address_id FROM users WHERE wechat_open_id = ?`
	var user model.User
	err := s.deps.DB.QueryRowContext(ctx, query, openID).Scan(
		&user.ID,
		&user.WeChatOpenID,
		&user.Nickname,
		&user.AvatarURL,
		&user.Phone,
		&user.DefaultAddrID,
	)
	switch {
	case err == nil:
		return &user, nil
	case errors.Is(err, sql.ErrNoRows):
		user = model.User{
			ID:           uid.New("usr_"),
			WeChatOpenID: openID,
			Nickname:     fmt.Sprintf("微信用户%s", openID),
			AvatarURL:    "https://example.com/avatar.png",
		}
		const insert = `INSERT INTO users (id, wechat_open_id, nickname, avatar_url) VALUES (?, ?, ?, ?)`
		if _, err := s.deps.DB.ExecContext(ctx, insert, user.ID, user.WeChatOpenID, user.Nickname, user.AvatarURL); err != nil {
			return nil, err
		}
		return &user, nil
	default:
		return nil, err
	}
}

func (s *userService) BindUser(ctx context.Context, user *model.User) error {
	if s.deps.DB == nil {
		return errUserDBUnavailable
	}
	if user == nil {
		return errors.New("user payload is nil")
	}
	if user.ID == "" {
		return errors.New("user id is required")
	}

	const stmt = `UPDATE users SET nickname = ?, avatar_url = ?, phone = ?, default_address_id = ?, updated_at = NOW() WHERE id = ?`
	res, err := s.deps.DB.ExecContext(ctx, stmt, user.Nickname, user.AvatarURL, user.Phone, user.DefaultAddrID, user.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("user %s not found", user.ID)
	}

	return nil
}

func (s *userService) ListAddresses(ctx context.Context, userID string) ([]model.Address, error) {
	if s.deps.DB == nil {
		return nil, errUserDBUnavailable
	}
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	const query = `SELECT id, user_id, recipient, phone, province, city, district, detail, postal_code, is_default FROM addresses WHERE user_id = ? ORDER BY updated_at DESC`
	rows, err := s.deps.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []model.Address
	for rows.Next() {
		var addr model.Address
		if err := rows.Scan(&addr.ID, &addr.UserID, &addr.Recipient, &addr.Phone, &addr.Province, &addr.City, &addr.District, &addr.Detail, &addr.PostalCode, &addr.IsDefault); err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (s *userService) CreateAddress(ctx context.Context, address *model.Address) (err error) {
	if s.deps.DB == nil {
		return errUserDBUnavailable
	}
	if address == nil {
		return errors.New("address payload is nil")
	}
	if address.UserID == "" {
		return errors.New("user id is required")
	}
	if address.ID == "" {
		address.ID = uid.New("addr_")
	}

	tx, err := s.deps.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	const insert = `INSERT INTO addresses (id, user_id, recipient, phone, province, city, district, detail, postal_code, is_default) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err = tx.ExecContext(ctx, insert, address.ID, address.UserID, address.Recipient, address.Phone, address.Province, address.City, address.District, address.Detail, address.PostalCode, address.IsDefault); err != nil {
		return err
	}

	if address.IsDefault {
		if _, err = tx.ExecContext(ctx, `UPDATE addresses SET is_default = FALSE WHERE user_id = ? AND id <> ?`, address.UserID, address.ID); err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `UPDATE users SET default_address_id = ? WHERE id = ?`, address.ID, address.UserID); err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (s *userService) UpdateAddress(ctx context.Context, address *model.Address) (err error) {
	if s.deps.DB == nil {
		return errUserDBUnavailable
	}
	if address == nil {
		return errors.New("address payload is nil")
	}
	if address.ID == "" {
		return errors.New("address id is required")
	}

	tx, err := s.deps.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var userID string
	if address.UserID != "" {
		userID = address.UserID
	} else {
		if err = tx.QueryRowContext(ctx, `SELECT user_id FROM addresses WHERE id = ?`, address.ID).Scan(&userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("address %s not found", address.ID)
			}
			return err
		}
	}

	const update = `UPDATE addresses SET recipient = ?, phone = ?, province = ?, city = ?, district = ?, detail = ?, postal_code = ?, is_default = ?, updated_at = NOW() WHERE id = ?`
	res, execErr := tx.ExecContext(ctx, update, address.Recipient, address.Phone, address.Province, address.City, address.District, address.Detail, address.PostalCode, address.IsDefault, address.ID)
	if execErr != nil {
		err = execErr
		return err
	}

	rowsAffected, execErr := res.RowsAffected()
	if execErr != nil {
		err = execErr
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("address %s not found", address.ID)
	}

	if address.IsDefault {
		if _, err = tx.ExecContext(ctx, `UPDATE addresses SET is_default = FALSE WHERE user_id = ? AND id <> ?`, userID, address.ID); err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `UPDATE users SET default_address_id = ? WHERE id = ?`, address.ID, userID); err != nil {
			return err
		}
	} else {
		if _, err = tx.ExecContext(ctx, `UPDATE users SET default_address_id = NULL WHERE id = ? AND default_address_id = ?`, userID, address.ID); err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (s *userService) DeleteAddress(ctx context.Context, addressID string) (err error) {
	if s.deps.DB == nil {
		return errUserDBUnavailable
	}
	if addressID == "" {
		return errors.New("address id is required")
	}

	tx, err := s.deps.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var userID string
	if err = tx.QueryRowContext(ctx, `SELECT user_id FROM addresses WHERE id = ?`, addressID).Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("address %s not found", addressID)
		}
		return err
	}

	res, execErr := tx.ExecContext(ctx, `DELETE FROM addresses WHERE id = ?`, addressID)
	if execErr != nil {
		err = execErr
		return err
	}

	affected, execErr := res.RowsAffected()
	if execErr != nil {
		err = execErr
		return err
	}
	if affected == 0 {
		return fmt.Errorf("address %s not found", addressID)
	}

	if _, err = tx.ExecContext(ctx, `UPDATE users SET default_address_id = NULL WHERE id = ? AND default_address_id = ?`, userID, addressID); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
