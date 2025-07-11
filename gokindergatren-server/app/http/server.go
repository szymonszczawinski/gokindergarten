package http

import (
	"context"
	"embed"
	"fmt"
	"gokindergarten/app/api"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type StaticContent struct {
	PublicDir embed.FS
}

type IHttpServer interface {
	AddHandler(c IViewHandler)
}

type httpServer struct {
	server     *http.Server
	router     *gin.Engine
	routes     *Routes
	fileSystem embed.FS
	group      *errgroup.Group
	ctx        context.Context
}

func NewHTTPServer(context context.Context, group *errgroup.Group, port int, staticContent embed.FS) *httpServer {
	ginRouter := createGinRouter(staticContent)
	routes := configureMainRoutes(ginRouter)
	instance := httpServer{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: ginRouter,
		},
		router:     ginRouter,
		routes:     routes,
		fileSystem: staticContent,
		group:      group,
		ctx:        context,
	}
	return &instance
}

func (s *httpServer) Start() {
	s.group.Go(func() error {
		s.group.Go(func() error {
			// service connections
			slog.Info("http server staring on", "address", s.server.Addr)
			if err := s.server.ListenAndServe(); err != nil {
				slog.Info("http server error", "error", err)
				return err
			}
			return nil
		})
		<-s.ctx.Done()
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
		// Listen for the interrupt signal.
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		} else {
			slog.Info("Server Shutdown OK")
		}
		// catching ctx.Done(). timeout of 5 seconds.
		// select {
		// case <-ctx.Done():
		// 	slog.Info("timeout of 5 seconds.")
		// }
		slog.Info("Server exiting")
		return nil
	})
}

func (s *httpServer) AddHandler(path string, vh IViewHandler) {
	vh.ConfigureRoutes(path, *s.routes)
}

func createGinRouter(fs embed.FS) *gin.Engine {
	engine := gin.Default()
	engine.StaticFS(api.AppPublicPath, http.FS(fs))
	// cookieOptions := sessions.Options{
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteLaxMode,
	// 	// Domain:   "gopt",
	// 	MaxAge: 60 * 5,
	// }
	cookieStore := cookie.NewStore([]byte(os.Getenv("SECRET")))
	// cookieStore.Options(cookieOptions)
	engine.Use(sessions.Sessions("mysession", cookieStore))
	return engine
}

func configureMainRoutes(router *gin.Engine) *Routes {
	rootRoute := router.Group(api.AppRootPath)
	apiRoute := rootRoute.Group(api.AppAPIPath)
	viewsRoute := rootRoute.Group(api.AppViewsPath)

	// apiRoute.Use(auth.SessionAuth)
	//	viewsRoute.Use(handlers.SessionAuth)
	routes := NewRoutes(rootRoute, viewsRoute, apiRoute)
	return routes
}
