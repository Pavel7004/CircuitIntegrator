package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/http/handler"
	"github.com/Pavel7004/GraphPlot/pkg/infra/config"
)

// @title           GraphPlot
// @version         0.3.0
// @description     This is a server that generates plots

// @contact.name   Kovalev Pavel
// @contact.email  kovalev5690@gmail.com

// @license.name   GPL-3.0
// @license.url    https://www.gnu.org/licenses/gpl-3.0.html

type Server struct {
	server    *http.Server
	router    *gin.Engine
	isRunning bool
	isDebug   bool

	handler *handler.Handler
}

func New() *Server {
	cfg := config.Get()
	server := new(Server)

	server.router = gin.New()
	server.server = &http.Server{
		Addr:           cfg.Hostname + ":" + cfg.Port,
		Handler:        server.router,
		ReadTimeout:    cfg.Timeout,
		WriteTimeout:   cfg.Timeout,
		MaxHeaderBytes: 100 * 1024 * 8, // 100 KiB
	}
	server.handler = handler.New()
	server.isDebug = true

	server.prepareRouter()

	return server
}

func EnableRelease() {
	gin.SetMode(gin.ReleaseMode)
}

func (s *Server) Run() error {
	s.isRunning = true
	return s.server.ListenAndServe()
}

func (s *Server) prepareRouter() {
	s.router.GET("/", s.handler.GetIndexPage)
	s.router.GET("/gate", s.handler.OpenWebsocketConnection)

	s.router.Static("/static", "./static")
}
