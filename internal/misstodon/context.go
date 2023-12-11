package misstodon

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/labstack/echo/v4"
)

type Ctx struct {
	m sync.Map
}

type Context interface {
	ProxyServer() string
	Token() *string
	UserID() *string
	HOST() *string
}

func ContextWithEchoContext(eCtx echo.Context, tokenRequired ...bool) (*Ctx, error) {
	c := &Ctx{}
	if server, ok := eCtx.Get("proxy-server").(string); ok {
		c.SetProxyServer(server)
	}
	token, _ := utils.GetHeaderToken(eCtx.Request().Header)
	tokenArr := strings.Split(token, ".")
	if len(tokenArr) >= 2 {
		c.SetUserID(tokenArr[0])
		c.SetToken(tokenArr[1])
	}
	if len(tokenRequired) > 0 && tokenRequired[0] {
		if token == "" {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "the access token is invalid")
		}
		if len(tokenArr) < 2 {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "the access token is invalid")
		}
	}
	c.SetHOST(eCtx.Request().Host)
	return c, nil
}

func ContextWithValues(proxyServer, token string) *Ctx {
	c := &Ctx{}
	c.SetProxyServer(proxyServer)
	c.SetToken(token)
	return c
}

func (*Ctx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*Ctx) Done() <-chan struct{} {
	return nil
}

func (*Ctx) Err() error {
	return nil
}

func (c *Ctx) Value(key any) any {
	if val, ok := c.m.Load(key); ok {
		return val
	}
	return nil
}

func (c *Ctx) SetValue(key any, val any) {
	c.m.Store(key, val)
}

func (c *Ctx) String(key string) *string {
	if val, ok := c.m.Load(key); ok {
		valStr := val.(string)
		return &valStr
	}
	return nil
}

func (c *Ctx) ProxyServer() string {
	return *c.String("proxy-server")
}

func (c *Ctx) SetProxyServer(val string) {
	c.SetValue("proxy-server", val)
}

func (c *Ctx) Token() *string {
	return c.String("token")
}

func (c *Ctx) SetToken(val string) {
	c.SetValue("token", val)
}

func (c *Ctx) UserID() *string {
	return c.String("user_id")
}

func (c *Ctx) SetUserID(val string) {
	c.SetValue("user_id", val)
}

func (c *Ctx) HOST() *string {
	return c.String("host")
}

func (c *Ctx) SetHOST(val string) {
	c.SetValue("host", val)
}
