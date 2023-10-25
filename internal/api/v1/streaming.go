package v1

import (
	"context"
	"net/http"

	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey/streaming"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var wsUpgrade = websocket.Upgrader{
	ReadBufferSize:  4096, // we don't expect reads
	WriteBufferSize: 4096,
	NegotiateSubprotocol: func(r *http.Request) (string, error) {
		return r.Header.Get("Sec-Websocket-Protocol"), nil
	},
	CheckOrigin: func(r *http.Request) bool { return true },
}

func StreamingRouter(e *echo.Group) {
	e.GET("/streaming", StreamingHandler)
}

func StreamingHandler(c echo.Context) error {
	var token string
	if token = c.QueryParam("access_token"); token == "" {
		if token = c.Request().Header.Get("Sec-Websocket-Protocol"); token == "" {
			return errors.New("no access token provided")
		}
	}
	server := c.Get("proxy-server").(string)

	conn, err := wsUpgrade.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan models.StreamEvent)
	defer close(ch)
	go func() {
		if err := streaming.Streaming(ctx, server, token, ch); err != nil {
			log.Debug().Caller().Err(err).Msg("Streaming error")
		}
		_ = conn.Close()
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-ch:
				log.Debug().Caller().Any("event", event).Msg("Streaming")
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		_, _, err = conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				cancel()
				return nil
			}
			return err
		}
	}
}
