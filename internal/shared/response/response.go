package response

import "github.com/gin-gonic/gin"

type Envelope struct {
	Data any `json:"data"`
}

type Error struct {
	Message string `json:"message"`
}

type Mutation struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

func OK(c *gin.Context, status int, data any) {
	c.JSON(status, Envelope{Data: data})
}

func Success(c *gin.Context, status int, message, id string) {
	c.JSON(status, Mutation{Message: message, ID: id})
}

func Fail(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, Error{Message: message})
}
