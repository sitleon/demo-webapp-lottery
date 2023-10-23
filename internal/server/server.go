package server

import (
	"demo-webapp-lottery/internal/app"
	"demo-webapp-lottery/internal/controller"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	*http.Server
	cfg *app.AppConfig

	controllers []controller.RestController
}

func NewServer(cfg *app.AppConfig, controllers []controller.RestController) *Server {
	r := gin.Default()
	srv := &Server{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Http.Server.Port),
			Handler: r,
		},
		cfg:         cfg,
		controllers: controllers,
	}

	srv.RegisterRoute(r)
	return srv
}

func (srv *Server) RegisterRoute(r *gin.Engine) {
	// health check
	r.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// register route
	for _, ctr := range srv.controllers {
		ctr.RegisterRoute(r.Group("/api"))
	}
}

func (srv *Server) Serve() {
	logrus.Infof("app server start to listen, %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
