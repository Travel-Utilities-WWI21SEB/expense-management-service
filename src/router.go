package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
	router := gin.Default()
	// apiv1 := router.Group("/api/v1")

	router.Handle(http.MethodGet, "/lifecheck", LifeCheckHandler())

	return router
}
