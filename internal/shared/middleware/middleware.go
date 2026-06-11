package middleware

import (
	"bookstore-api/internal/platform/logging"
	"bookstore-api/internal/shared/response"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"
const ErrorMessageKey = "error_message"

type ErrorMapper func(err error) (status int, message string, handled bool)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = newRequestID()
		}

		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func ErrorHandler(mappers ...ErrorMapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() || len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		for _, mapper := range mappers {
			if status, message, handled := mapper(err); handled {
				c.Errors = nil
				response.Fail(c, status, message)
				return
			}
		}

		message := c.GetString(ErrorMessageKey)
		if message == "" {
			message = "internal server error"
		}

		logging.Logger.Error().
			Err(err).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("client_ip", c.ClientIP()).
			Str("request_id", c.GetString(RequestIDKey)).
			Msg(message)

		c.Errors = nil
		response.Fail(c, http.StatusInternalServerError, message)
	}
}

func AddError(c *gin.Context, err error, message string) {
	c.Set(ErrorMessageKey, message)
	_ = c.Error(err)
}

func newRequestID() string {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(value)
}
