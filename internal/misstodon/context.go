package misstodon

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/labstack/echo/v4"
)

type Context struct {
	m sync.Map
}

func ContextWithEchoContext(eCtx echo.Context, tokenRequired ...bool) (*Context, error) {
	c := &Context{}
	if server, ok := eCtx.Get("server").(string); ok {
		c.SetServer(server)
	}
	if len(tokenRequired) > 0 && tokenRequired[0] {
		token, err := utils.GetHeaderToken(eCtx.Request().Header)
		if err != nil && (len(tokenRequired) > 0 && tokenRequired[0]) {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "the access token is invalid")
		}
		arr := strings.Split(token, ".")
		if len(arr) < 2 {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "the access token is invalid")
		}
		c.SetUserID(arr[0])
		c.SetToken(arr[1])
	}
	return c, nil
}

func ContextWithValues(server, token string) *Context {
	c := &Context{}
	c.SetServer(server)
	c.SetToken(token)
	return c
}

func (*Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*Context) Done() <-chan struct{} {
	return nil
}

func (*Context) Err() error {
	return nil
}

func (c *Context) Value(key any) any {
	if val, ok := c.m.Load(key); ok {
		return val
	}
	return nil
}

func (c *Context) SetValue(key any, val any) {
	c.m.Store(key, val)
}

func (c *Context) String(key string) *string {
	if val, ok := c.m.Load(key); ok {
		valStr := val.(string)
		return &valStr
	}
	return nil
}

func (c *Context) Server() string {
	return *c.String("server")
}

func (c *Context) SetServer(val string) {
	c.SetValue("server", val)
}

func (c *Context) Token() *string {
	return c.String("token")
}

func (c *Context) SetToken(val string) {
	c.SetValue("token", val)
}

func (c *Context) UserID() *string {
	return c.String("user_id")
}

func (c *Context) SetUserID(val string) {
	c.SetValue("user_id", val)
}
