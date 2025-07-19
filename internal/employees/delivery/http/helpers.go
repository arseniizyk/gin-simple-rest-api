package delivery

import (
	"log"

	"github.com/gin-gonic/gin"
)

func sendError(c *gin.Context, status int, msg string, err error) {
	log.Println(msg+":", err)
	c.JSON(status, ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	})
}
