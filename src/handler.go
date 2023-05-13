package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/model"
)

func LifeCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := model.LifeCheckResponse{
			Alive:   true,
			Version: "1.0.0",
		}

		c.JSON(http.StatusOK, response)
	}
}
