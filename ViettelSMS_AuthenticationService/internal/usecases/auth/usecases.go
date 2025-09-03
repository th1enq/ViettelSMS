package auth

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/golang-jwt/jwt"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/dto"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/entity"
	domain "github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/errors"
	repo "github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/repository"
	srv "github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/service"
	rdb "github.com/th1enq/ViettelSMS_AuthenticationService/internal/infrastucture/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type usecase struct {
	repo      repo.Repository
	cache     rdb.CacheEngine
	password  srv.PasswordService
	jwtSecret string
	logger    *zap.Logger
}

func NewUseCase(
	repo repo.Repository,
	password srv.PasswordService,
	cache rdb.CacheEngine,
	jwtSecret string,
	logger *zap.Logger,
) UseCase {
	return &usecase{
		repo:      repo,
		password:  password,
		cache:     cache,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// decodePayload decodes a map[string]interface{} payload into a struct using JSON tags
func (u *usecase) decodePayload(payload map[string]interface{}, result interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  result,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		u.logger.Error("Failed to create decoder", zap.Error(err))
		return err
	}

	if err := decoder.Decode(payload); err != nil {
		u.logger.Error("Failed to decode payload", zap.Error(err), zap.Any("payload", payload))
		return err
	}

	return nil
}

func (u *usecase) Login(ctx context.Context, userName, password string) (*dto.LoginResponse, error) {
	u.logger.Info("Login attempt", zap.String("username", userName))
	user, err := u.repo.GetUserByUsername(ctx, userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.logger.Warn("User not found", zap.String("username", userName))
			return nil, domain.ErrUserNotFound
		} else {
			u.logger.Error("Failed to retrieve user", zap.String("username", userName), zap.Error(err))
			return nil, domain.ErrInternalServer
		}
	}

	ok, err := u.password.Verify(user.Password, password)
	if err != nil {
		u.logger.Error("Failed to verify password", zap.String("username", userName), zap.Error(err))
		return nil, domain.ErrInternalServer
	}
	if !ok {
		u.logger.Warn("Invalid credentials", zap.String("username", userName))
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		u.logger.Error("Failed to generate access token", zap.String("username", userName), zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	refreshToken, err := u.generateRefreshToken(user)
	if err != nil {
		u.logger.Error("Failed to generate refresh token", zap.String("username", userName), zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	u.logger.Info("Login successful", zap.String("username", userName))
	return &dto.LoginResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (u *usecase) RefreshToken(ctx context.Context, userID uint) (*string, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		u.logger.Error("Failed to retrieve user", zap.Uint("userID", userID), zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		u.logger.Error("Failed to generate access token", zap.Uint("userID", userID), zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	u.logger.Info("Token refreshed successfully", zap.Uint("userID", userID))
	return accessToken, nil
}

func (u *usecase) CreateAuthUser(ctx context.Context, payload map[string]interface{}) error {
	u.logger.Info("Creating user", zap.Any("payload", payload))

	var req dto.UserCreate
	if err := u.decodePayload(payload, &req); err != nil {
		return err
	}

	u.logger.Info("Decoded user create request", zap.String("username", req.Username))

	user := &entity.AuthUser{
		Username: req.Username,
		Password: req.Password,
		Blocked:  req.Blocked,
		Scopes:   req.Scopes,
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	u.logger.Info("User created successfully", zap.String("username", req.Username))

	return nil
}

func (u *usecase) UpdateAuthUser(ctx context.Context, payload map[string]interface{}) error {
	u.logger.Info("Updating user", zap.Any("payload", payload))

	var req dto.UserUpdate
	if err := u.decodePayload(payload, &req); err != nil {
		return err
	}

	u.logger.Info("Decoded user update request", zap.String("user_name", req.UserName))

	user, err := u.getUserByUserName(ctx, req.UserName)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil
		}
		return err
	}

	user.Blocked = req.Blocked

	if err := u.repo.UpdateUser(ctx, user); err != nil {
		return err
	}

	u.logger.Info("User updated successfully", zap.String("user_name", req.UserName))

	return nil
}

func (u *usecase) DeleteAuthUser(ctx context.Context, payload map[string]interface{}) error {
	u.logger.Info("Deleting user", zap.Any("payload", payload))

	var req dto.UserDelete
	if err := u.decodePayload(payload, &req); err != nil {
		return err
	}

	u.logger.Info("Decoded user delete request", zap.String("user_name", req.UserName))

	user, err := u.getUserByUserName(ctx, req.UserName)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil
		}
		return err
	}

	if err := u.repo.DeleteUser(ctx, user.ID); err != nil {
		u.logger.Error("Failed to delete user", zap.String("user_name", req.UserName), zap.Error(err))
		return domain.ErrInternalServer
	}

	u.logger.Info("User deleted successfully", zap.String("user_name", req.UserName))

	return nil
}

func (u *usecase) UpdateAuthPassword(ctx context.Context, payload map[string]interface{}) error {
	u.logger.Info("Updating user password", zap.Any("payload", payload))

	var req dto.UserPasswordUpdate
	if err := u.decodePayload(payload, &req); err != nil {
		return err
	}

	u.logger.Info("Decoded password update request", zap.String("user_name", req.UserName))

	user, err := u.getUserByUserName(ctx, req.UserName)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil
		}
		return err
	}

	user.Password = req.Password

	if err := u.repo.UpdateUser(ctx, user); err != nil {
		u.logger.Error("Failed to update user", zap.String("user_name", req.UserName), zap.Error(err))
		return domain.ErrInternalServer
	}

	u.logger.Info("User password updated successfully", zap.String("user_name", req.UserName))

	return nil
}

func (u *usecase) AddAuthScope(ctx context.Context, payload map[string]interface{}) error {
	u.logger.Info("Adding user scope", zap.Any("payload", payload))

	var req dto.UserScope
	if err := u.decodePayload(payload, &req); err != nil {
		return err
	}

	u.logger.Info("Decoded add scope request", zap.String("user_name", req.UserName), zap.String("scope", req.Scope))

	user, err := u.getUserByUserName(ctx, req.UserName)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil
		}
		return err
	}

	user.Scopes = append(user.Scopes, req.Scope)

	if err := u.repo.UpdateUser(ctx, user); err != nil {
		return err
	}

	u.logger.Info("User scope added successfully", zap.String("user_name", req.UserName), zap.String("scope", req.Scope))

	return nil
}

func (u *usecase) DeleteAuthScope(ctx context.Context, payload map[string]interface{}) error {
	u.logger.Info("Deleting user scope", zap.Any("payload", payload))

	var req dto.UserScope
	if err := u.decodePayload(payload, &req); err != nil {
		return err
	}

	u.logger.Info("Decoded delete scope request", zap.String("user_name", req.UserName), zap.String("scope", req.Scope))

	user, err := u.getUserByUserName(ctx, req.UserName)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil
		}
		return err
	}

	for i, scope := range user.Scopes {
		if scope == req.Scope {
			user.Scopes = slices.Delete(user.Scopes, i, i+1)
			break
		}
	}

	u.logger.Info("Current user scopes", zap.String("user_name", req.UserName), zap.Strings("scopes", user.Scopes))

	if err := u.repo.UpdateUser(ctx, user); err != nil {
		return err
	}

	u.logger.Info("User scope deleted successfully", zap.String("user_name", req.UserName), zap.String("scope", req.Scope))

	return nil
}

func (u *usecase) generateAccessToken(user *entity.AuthUser) (*string, error) {
	scopes := make([]string, len(user.Scopes))
	copy(scopes, user.Scopes)

	claims := jwt.MapClaims{
		"sub":     user.ID,
		"blocked": user.Blocked,
		"scopes":  scopes,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}
	return &signedToken, nil
}

func (u *usecase) generateRefreshToken(user *entity.AuthUser) (*string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}
	return &signedToken, nil
}

func (u *usecase) getUserByUserName(ctx context.Context, userName string) (*entity.AuthUser, error) {
	user, err := u.repo.GetUserByUsername(ctx, userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.logger.Warn("User not found", zap.String("user_name", userName))
			return nil, domain.ErrUserNotFound
		}
		u.logger.Error("Failed to retrieve user", zap.String("user_name", userName), zap.Error(err))
		return nil, domain.ErrInternalServer
	}
	return user, nil
}
