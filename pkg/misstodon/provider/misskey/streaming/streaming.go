package streaming

import (
	"context"
	"fmt"
	"time"

	"github.com/gizmo-ds/misstodon/internal/misstodon"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/pkg/misstodon/models"
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

func Streaming(ctx context.Context, mCtx misstodon.Context, token string, ch chan<- models.StreamEvent) error {
	u := fmt.Sprintf("wss://%s/streaming?i=%s&_t=%d", mCtx.ProxyServer(), token, time.Now().Unix())
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	_ = conn.WriteJSON(utils.Map{
		"type": "connect",
		"body": utils.Map{
			"channel": "main",
			"id":      xid.New().String(),
		},
	})

	done := false
	go func() {
		select {
		case <-ctx.Done():
			done = true
			_ = conn.Close()
		}
	}()

	for {
		var v models.MkStreamMessage
		if err = conn.ReadJSON(&v); err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				return nil
			}
			if done {
				return nil
			}
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		if done {
			return nil
		}
		ch <- v.ToStreamEvent()
	}
}
