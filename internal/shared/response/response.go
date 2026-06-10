package response

import "github.com/gin-gonic/gin"

type Envelope struct {
	Data any `json:"data"`
}

type Error struct {
	Message string `json:"message"`
}

func OK(c *gin.Context, status int, data any) {
	c.JSON(status, Envelope{Data: data})
}

func Fail(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, Error{Message: message})
}
