package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/config"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/delivery/http/controller"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/delivery/http/middleware"
	"go.uber.org/zap"
)

type (
	Server interface {
		Start(ctx context.Context) error
	}

	server struct {
		config        *config.Config
		controller    *controller.Controller
		jwtMiddleware middleware.JWTMiddleware
		logger        *zap.Logger
	}
)

func NewHttpServer(
	config *config.Config,
	controller *controller.Controller,
	jwtMiddleware middleware.JWTMiddleware,
	logger *zap.Logger,
) Server {
	return &server{
		config:        config,
		controller:    controller,
		jwtMiddleware: jwtMiddleware,
		logger:        logger,
	}
}

func (s *server) RegisterRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", s.controller.Login)
		auth.POST("/refresh/{user_id}", s.controller.RefreshToken)
	}

	return router
}

func (s *server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler: s.RegisterRoutes(),
	}

	s.logger.Info("Starting HTTP server", zap.String("address", server.Addr))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal("Failed to start HTTP server", zap.Error(err))
	}
	return nil
}
