package handlers

import (
	"net/http"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/gin-gonic/gin"
)

func LifeCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := &models.LifeCheckResponse{
			Alive:   true,
			Version: "1.0.0",
		}

		c.JSON(http.StatusOK, response)
	}
}
