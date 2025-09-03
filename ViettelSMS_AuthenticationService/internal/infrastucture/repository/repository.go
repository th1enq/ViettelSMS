package repository

import (
	"context"

	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/entity"
	repo "github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/repository"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/infrastucture/postgres"
)

type repository struct {
	db postgres.DBEngine
}

func NewRepository(
	db postgres.DBEngine,
) repo.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetUserByID(ctx context.Context, id uint) (*entity.AuthUser, error) {
	var user entity.AuthUser
	if err := r.db.GetDB().WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*entity.AuthUser, error) {
	var user entity.AuthUser
	if err := r.db.GetDB().WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) CreateUser(ctx context.Context, user *entity.AuthUser) error {
	if err := r.db.GetDB().WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateUser(ctx context.Context, user *entity.AuthUser) error {
	if err := r.db.GetDB().WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteUser(ctx context.Context, id uint) error {
	if err := r.db.GetDB().WithContext(ctx).Delete(&entity.AuthUser{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
