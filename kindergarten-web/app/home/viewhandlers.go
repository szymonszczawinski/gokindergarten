// Package home
package home

import (
	"kindergarten-web/app/http"
	"kindergarten-web/views/home"

	"github.com/gin-gonic/gin"
)

type homeHandler struct{}

func NewHomeHandler() homeHandler {
	instance := homeHandler{}
	return instance
}

func (handler homeHandler) ConfigureRoutes(path string, routes http.Routes) {
	routes.Root().GET("/", handler.homePage)
	routes.Root().GET("/error", handler.errorPage)
}

func (handler *homeHandler) homePage(c *gin.Context) {
	isHxRequest := c.GetHeader("HX-Request")
	if isHxRequest == "true" {
		// home.Home(true).Render(c.Request.Context(), c.Writer)
		c.JSON(200, gin.H{"Hello": "World"})
	} else {
		// c.JSON(200, gin.H{"Hello": "World"})
		home.Home(false).Render(c.Request.Context(), c.Writer)
	}
}

func (handler *homeHandler) errorPage(c *gin.Context) {
	// 	errors.Error("SUPER ERROR").Render(c.Request.Context(), c.Writer)
}
