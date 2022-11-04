package handler

import (
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	session "github.com/Pavel7004/GraphPlot/pkg/components/web-session"
	"github.com/Pavel7004/GraphPlot/pkg/domain"
)

type Handler struct {
	upgrader *websocket.Upgrader

	sessions []*session.Session
}

func New() *Handler {
	return &Handler{
		upgrader: &websocket.Upgrader{HandshakeTimeout: 5 * time.Second, ReadBufferSize: 1024 * 1024, WriteBufferSize: 1024 * 1024},
		sessions: []*session.Session{},
	}
}

func (h *Handler) SendError(c *gin.Context, err error) {
	if e, ok := err.(*domain.Error); ok {
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
		log.Warn().Err(err).Msg("Failed to parse \"templates/index.html\"")
	}

	if err := tmpl.Execute(c.Writer, struct {
		Integrators []string
	}{
		domain.IntegratorsNames,
	}); err != nil {
		log.Error().Err(err).Msg("Failed to execute \"templates/index.html\"")
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
