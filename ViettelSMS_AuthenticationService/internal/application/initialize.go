package application

import (
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/config"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/delivery/consumer"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/delivery/http"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/delivery/http/controller"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/delivery/http/middleware"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/delivery/http/presenter"
	password "github.com/th1enq/ViettelSMS_AuthenticationService/internal/infrastucture/bcrypt"
	consumerGroup "github.com/th1enq/ViettelSMS_AuthenticationService/internal/infrastucture/kafka/consumer"
	log "github.com/th1enq/ViettelSMS_AuthenticationService/internal/infrastucture/logger"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/infrastucture/postgres"
	rdb "github.com/th1enq/ViettelSMS_AuthenticationService/internal/infrastucture/redis"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/infrastucture/repository"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/usecases/auth"
)

func InitApp() (*Application, error) {
	config := config.LoadConfig()

	logger, err := log.LoadLogger(config)
	if err != nil {
		return nil, err
	}

	db, err := postgres.NewPostgresDB(config, logger)
	if err != nil {
		return nil, err
	}

	rdb, err := rdb.NewRedisDB(config)
	if err != nil {
		return nil, err
	}

	repo := repository.NewRepository(db)

	passwordSrv := password.NewBcryptService()
	usecase := auth.NewUseCase(
		repo,
		passwordSrv,
		rdb,
		config.JWT.Secret,
		logger,
	)

	presenter := presenter.NewPresenter()
	middleware := middleware.NewJWTMiddleware(presenter, []byte(config.JWT.Secret))
	controller := controller.NewController(logger, usecase, presenter)

	httpServer := http.NewHttpServer(config, controller, middleware, logger)

	userConsumer, err := consumerGroup.NewConsumer(
		config,
		logger,
		config.Consumer.UserAuth,
	)

	userHandler := consumer.NewUserHandler(logger, usecase)

	if err != nil {
		return nil, err
	}

	rootConsumer := consumer.NewRoot(
		logger,
		userConsumer,
		userHandler,
	)

	app := NewApplication(httpServer, rootConsumer, logger)
	return app, nil
}
