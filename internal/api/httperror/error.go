package httperror

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

type ServerError struct {
	TraceID string `json:"trace_id,omitempty"`
	Error   string `json:"error"`
}

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if err == nil {
		return
	}

	info := ServerError{Error: err.Error()}
	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		info.Error = fmt.Sprint(he.Message)
	}

	if code == http.StatusInternalServerError {
		errorID := xid.New().String()
		info = ServerError{
			TraceID: errorID,
			Error:   "Internal Server Error",
		}
		log.Warn().Err(err).
			Str("user_agent", c.Request().UserAgent()).
			Str("trace_id", errorID).
			Int("code", code).
			Msg("Server Error")
	}
	_ = c.JSON(code, info)
}
