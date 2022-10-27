package http

import (
	"net/http"
	"time"

	"github.com/Pavel7004/GraphPlot/pkg/infra/config"
	"github.com/gin-gonic/gin"
)

// @title           GraphPlot
// @version         0.1
// @description     This is a server that generates plots

// @contact.name   Kovalev Pavel
// @contact.email  kovalev5690@gmail.com

// @license.name   GPL-3.0
// @license.url    https://www.gnu.org/licenses/gpl-3.0.html

type Server struct {
	server    *http.Server
	router    *gin.Engine
	isRunning bool
}

func New() *Server {
	server := new(Server)
	cfg := config.Get()

	server.router = gin.New()
	server.server = &http.Server{
		Addr:           cfg.Hostname + ":" + cfg.Port,
		Handler:        server.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 100 * 1024 * 8, // 100 KiB
	}

	server.prepareRouter()

	return server
}

func (s *Server) Run() error {
	s.isRunning = true
	return s.server.ListenAndServe()
}

func (s *Server) prepareRouter() {}
