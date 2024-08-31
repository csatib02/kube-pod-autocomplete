package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

// HandleHTTPError is a utility function to handle HTTP errors in Gin handlers.
func HandleHTTPError(c *gin.Context, err error) {
	httpErr := &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}

	c.JSON(httpErr.Code, httpErr)
}
