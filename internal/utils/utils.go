package utils

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseID(c *gin.Params) (int, error) {
	idStr := c.ByName("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Println("cant parse id to int", err.Error())
		return 0, err
	}

	return id, nil
}
