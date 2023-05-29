package misskey

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/go-resty/resty/v2"
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

func isucceed(resp *resty.Response, statusCode int, codes ...string) error {
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusNoContent, statusCode:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrUnauthorized
	case http.StatusNotFound:
		return ErrNotFound
	}

	var result struct {
		Error struct {
			Code string `json:"code"`
			Msg  string `json:"message"`
		} `json:"error"`
	}
	body := resp.Body()
	if body != nil {
		err := json.Unmarshal(body, &result)
		if err != nil {
			return ServerError{Code: 500, Message: err.Error()}
		}
	}
	if utils.Contains(codes, result.Error.Code) {
		return nil
	}
	return errors.New(result.Error.Msg)
}
