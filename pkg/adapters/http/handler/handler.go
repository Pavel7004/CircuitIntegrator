package handler

import (
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	session "github.com/Pavel7004/GraphPlot/pkg/components/web-session"
	"github.com/Pavel7004/GraphPlot/pkg/domain"
	"github.com/Pavel7004/GraphPlot/pkg/infra/config"
)

type Handler struct {
	upgrader *websocket.Upgrader

	sessions []*session.Session
}

func New() *Handler {
	cfg := config.Get()
	h := new(Handler)

	h.upgrader = &websocket.Upgrader{
		HandshakeTimeout: cfg.Timeout,
		ReadBufferSize:   1024 * 1024,
		WriteBufferSize:  1024 * 1024,
	}
	h.sessions = []*session.Session{}

	return h
}

func (h *Handler) SendError(c *gin.Context, err error) {
	if e, ok := err.(*domain.Error); ok { //nolint
		c.JSON(e.CodeHTTP, e)
	} else {
		c.JSON(500, &domain.Error{
			Code:    "unknown_error",
			Message: err.Error(),
		})
	}
}

func (h *Handler) GetIndexPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("[ERROR] Failed to parse \"templates/index.html\". err = %v", err)
	}

	if err := tmpl.Execute(c.Writer, struct {
		Integrators []string
	}{
		domain.IntegratorsNames,
	}); err != nil {
		log.Printf("[ERROR] to execute \"templates/index.html\". err = %v", err)
	}
}

func (h *Handler) OpenWebsocketConnection(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.SendError(c, err)
		return
	}

	h.sessions = append(h.sessions, session.New(conn))
}
