package consumer

import (
	"context"

	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/infrastucture/kafka/consumer"
	"go.uber.org/zap"
)

const (
	USER_TOPIC = "user-events"
)

type (
	Root interface {
		Start(ctx context.Context) error
	}

	root struct {
		logger *zap.Logger

		userConsumer consumer.Consumer
		userHandler  UserHandler
	}
)

func NewRoot(
	logger *zap.Logger,
	userConsumer consumer.Consumer,
	userHandler UserHandler,
) Root {
	return &root{
		logger:       logger,
		userConsumer: userConsumer,
		userHandler:  userHandler,
	}
}

func (r *root) Start(ctx context.Context) error {
	r.logger.Info("Starting Kafka consumer...")

	r.userConsumer.RegisterHandler(
		USER_TOPIC,
		func(ctx context.Context, queueName string, payload []byte) error {
			return r.userHandler.Handle(ctx, queueName, payload)
		},
	)

	r.logger.Info("Kafka consumer started, waiting for messages...")

	go func() {
		if err := r.userConsumer.Start(ctx); err != nil {
			r.logger.Error("Failed to start Kafka consumer", zap.Error(err))
		}
	}()
	return nil
}
