package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/delivery/http/presenter"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/dto"
	domain "github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/errors"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/usecases/auth"
	"go.uber.org/zap"
)

type Controller struct {
	logger    *zap.Logger
	usecase   auth.UseCase
	presenter presenter.Presenter
}

func NewController(
	logger *zap.Logger,
	usecase auth.UseCase,
	presenter presenter.Presenter,
) *Controller {
	return &Controller{
		logger:    logger,
		usecase:   usecase,
		presenter: presenter,
	}
}

// Login godoc
// @Summary Login
// @Description User login
// @Tags auth
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "Login request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	c.logger.Info("User logged in")

	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Failed to bind login request", zap.Error(err))
		c.presenter.InvalidRequest(ctx, "Invalid request", err)
		return
	}

	token, err := c.usecase.Login(ctx.Request.Context(), req.UserName, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.logger.Warn("User not found", zap.String("username", req.UserName))
			c.presenter.InvalidRequest(ctx, "Invalid request", err)
		} else {
			c.logger.Error("Failed to login", zap.Error(err))
			c.presenter.InternalError(ctx, "Internal server error", err)
		}
		return
	}

	c.presenter.LoginSuccess(ctx, "Login successful", token)
	c.logger.Info("User logged in successfully", zap.String("username", req.UserName))
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Refresh user token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/refresh [post]
func (c *Controller) RefreshToken(ctx *gin.Context) {
	c.logger.Info("Refreshing token")

	userID := ctx.GetUint("user_id")

	token, err := c.usecase.RefreshToken(ctx.Request.Context(), userID)
	if err != nil {
		c.logger.Error("Failed to refresh token", zap.Error(err))
		c.presenter.InternalError(ctx, "Internal server error", err)
		return
	}

	c.presenter.LoginSuccess(ctx, "Token refreshed successfully", *token)
	c.logger.Info("Token refreshed successfully")
}
