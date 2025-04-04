package response

import (
	"net/http"

	"github.com/datpham/user-service-ms/internal/errors"
	"github.com/gin-gonic/gin"
)

const (
	CREATED = "Created"
	OK      = "Ok"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Error      string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func NewResponse(statusCode int, err string, data interface{}) *Response {
	return &Response{
		StatusCode: statusCode,
		Error:      err,
		Data:       data,
	}
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, NewResponse(http.StatusCreated, "", data))
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, NewResponse(http.StatusOK, "", data))
}

func Error(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, NewResponse(statusCode, err.Error(), nil))
}

func ErrorWithData(c *gin.Context, statusCode int, err error, data interface{}) {
	c.JSON(statusCode, NewResponse(statusCode, err.Error(), data))
}

func ErrorService(c *gin.Context, err error) {
	customErr, ok := err.(*errors.CustomError)
	if !ok {
		customErr = errors.NewCustomError(errors.ErrInternalServer, err.Error())
	}

	statusCode := http.StatusInternalServerError
	switch customErr.Code {
	case errors.ErrInvalidRequest:
		statusCode = http.StatusBadRequest
	case errors.ErrInternalServer:
		statusCode = http.StatusInternalServerError
	case errors.ErrNotFound:
		statusCode = http.StatusNotFound
	case errors.ErrUnauthorized:
		statusCode = http.StatusUnauthorized
	case errors.ErrForbidden:
		statusCode = http.StatusForbidden
	case errors.ErrConflict:
		statusCode = http.StatusConflict
	}

	c.JSON(statusCode, NewResponse(statusCode, customErr.Error(), nil))
}

func Redirect(c *gin.Context, url string) {
	c.Redirect(http.StatusTemporaryRedirect, url)
}
