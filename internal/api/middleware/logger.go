package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

var Logger = middleware.RequestLoggerWithConfig(
	middleware.RequestLoggerConfig{
		LogMethod:  true,
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("proxy-server", c.Get("proxy-server").(string)).
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Str("latency", v.Latency.String()).
				Msg("Request")
			return nil
		},
	})
