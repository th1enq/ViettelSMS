package presenter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/response"
)

type Presenter interface {
	InvalidRequest(c *gin.Context, message string, err error)
	InternalError(c *gin.Context, message string, err error)

	LoginSuccess(c *gin.Context, message string, data interface{})
	Unauthorized(c *gin.Context, message string, err error)
	Forbidden(c *gin.Context, message string, err error)
}

type presenter struct{}

func NewPresenter() Presenter {
	return &presenter{}
}

func (p *presenter) InvalidRequest(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, response.NewErrorResponse(
		response.CodeBadRequest,
		message,
		err.Error(),
	))
}

func (p *presenter) InternalError(c *gin.Context, message string, err error) {
	c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
		response.CodeInternalServerError,
		message,
		err.Error(),
	))
}

func (p *presenter) LoginSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, response.NewSuccessResponse(
		response.CodeSuccess,
		message,
		data,
	))
}

func (p *presenter) Unauthorized(c *gin.Context, message string, err error) {
	c.JSON(http.StatusUnauthorized, response.NewErrorResponse(
		response.CodeUnauthorized,
		message,
		err.Error(),
	))
}

func (p *presenter) Forbidden(c *gin.Context, message string, err error) {
	c.JSON(http.StatusForbidden, response.NewErrorResponse(
		response.CodeForbidden,
		message,
		err.Error(),
	))
}
