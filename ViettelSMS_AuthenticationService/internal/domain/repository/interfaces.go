package repo

import (
	"context"

	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/entity"
)

type Repository interface {
	GetUserByUsername(ctx context.Context, username string) (*entity.AuthUser, error)
	GetUserByID(ctx context.Context, id uint) (*entity.AuthUser, error)

	CreateUser(ctx context.Context, user *entity.AuthUser) error
	UpdateUser(ctx context.Context, user *entity.AuthUser) error
	DeleteUser(ctx context.Context, id uint) error
}
