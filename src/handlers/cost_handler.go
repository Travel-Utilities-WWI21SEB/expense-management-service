package handlers

import (
	"net/http"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/controllers"
	"github.com/gin-gonic/gin"
)

func CreateCostEntryHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := costCtl.CreateCostEntry(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func UpdateCostEntryHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := costCtl.UpdateCostEntry(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetCostDetailsHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		response, err := costCtl.GetCostDetails(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteCostEntryHandler(costCtl controllers.CostCtl) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TO-DO
		ctx := c.Request.Context()

		err := costCtl.DeleteCostEntry(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}
