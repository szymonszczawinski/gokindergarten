// Package http
package http

import (
	"github.com/gin-gonic/gin"
)

// Interface for all handlers
type IViewHandler interface {
	// Perform(a Action) gin.HandlerFunc
	ConfigureRoutes(path string, routes Routes)
}

// Container for application routes
type Routes struct {
	root  *gin.RouterGroup
	views *gin.RouterGroup
	apis  *gin.RouterGroup
}

func NewRoutes(root, views, apis *gin.RouterGroup) *Routes {
	return &Routes{
		root:  root,
		views: views,
		apis:  apis,
	}
}

// Get main application route: /gopt
func (r Routes) Root() *gin.RouterGroup {
	return r.root
}

// Get application views route: /gopt/views
func (r Routes) Views() *gin.RouterGroup {
	return r.views
}

// Get application API route: /gopt/api
func (r Routes) Apis() *gin.RouterGroup {
	return r.apis
}
