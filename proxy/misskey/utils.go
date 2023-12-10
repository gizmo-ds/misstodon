package misskey

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/gizmo-ds/misstodon/pkg/httpclient"
)

var (
	ErrUnauthorized = errors.New("invalid token")
)

type ServerError struct {
	Code    int
	Message string
}

func (e ServerError) Error() string {
	return e.Message
}

func isucceed(resp httpclient.Response, statusCode int, codes ...string) error {
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent, statusCode:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrUnauthorized
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusTooManyRequests:
		return ErrRateLimit
	}

	var result struct {
		Error struct {
			Code string `json:"code"`
			Msg  string `json:"message"`
		} `json:"error"`
	}
	if strings.Contains(resp.Header().Get("Content-Type"), "application/json") {
		body := resp.Body()
		if body != nil {
			err := json.Unmarshal(body, &result)
			if err != nil {
				return ServerError{Code: 500, Message: err.Error()}
			}
		}
	}
	if utils.Contains(codes, result.Error.Code) {
		return nil
	}
	return errors.New(result.Error.Msg)
}

func makeBody(ctx Context, m utils.Map) utils.Map {
	r := utils.Map{}
	token := ctx.Token()
	if token != nil && *token != "" {
		r["i"] = token
	}
	for k, v := range m {
		r[k] = v
	}
	return r
}
