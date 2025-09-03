package consumer

import (
	"context"
	"encoding/json"

	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/dto"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/usecases/auth"
	"go.uber.org/zap"
)

type UserHandler interface {
	Handle(ctx context.Context, topic string, payload []byte) error
}

type userHandler struct {
	logger  *zap.Logger
	usecase auth.UseCase
}

func NewUserHandler(
	logger *zap.Logger,
	usecase auth.UseCase,
) UserHandler {
	return &userHandler{
		logger:  logger,
		usecase: usecase,
	}
}

func (u *userHandler) Handle(ctx context.Context, topic string, payload []byte) error {
	u.logger.Info("Handling user event message", zap.String("topic", topic), zap.ByteString("payload", payload))

	var msg dto.UserEvent
	if err := json.Unmarshal(payload, &msg); err != nil {
		u.logger.Error("failed to unmarshal user event", zap.Error(err))
		return err
	}

	switch msg.Event {
	case "user.created":
		return u.usecase.CreateAuthUser(ctx, msg.Payload)
	case "user.updated":
		return u.usecase.UpdateAuthUser(ctx, msg.Payload)
	case "user.deleted":
		return u.usecase.DeleteAuthUser(ctx, msg.Payload)
	case "user.updated_password":
		return u.usecase.UpdateAuthPassword(ctx, msg.Payload)
	case "user.added_scope":
		return u.usecase.AddAuthScope(ctx, msg.Payload)
	case "user.deleted_scope":
		return u.usecase.DeleteAuthScope(ctx, msg.Payload)
	default:
		u.logger.Warn("unknown user event", zap.String("event", msg.Event))
		return nil
	}
}
