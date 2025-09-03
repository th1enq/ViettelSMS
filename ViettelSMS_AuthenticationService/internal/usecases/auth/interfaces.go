package auth

import (
	"context"

	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/dto"
)

type UseCase interface {
	Login(ctx context.Context, userName, password string) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, userID uint) (*string, error)

	CreateAuthUser(ctx context.Context, payload map[string]interface{}) error
	UpdateAuthUser(ctx context.Context, payload map[string]interface{}) error
	DeleteAuthUser(ctx context.Context, payload map[string]interface{}) error
	UpdateAuthPassword(ctx context.Context, payload map[string]interface{}) error
	AddAuthScope(ctx context.Context, payload map[string]interface{}) error
	DeleteAuthScope(ctx context.Context, payload map[string]interface{}) error
}
