package v1

import (
	"bytes"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func MissingImageHandler(c echo.Context) error {
	png1x1 := []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13,
		73, 72, 68, 82, 0, 0, 0, 1, 0, 0, 0, 1, 8, 4, 0, 0, 0, 181,
		28, 12, 2, 0, 0, 0, 11, 73, 68, 65, 84, 120, 218, 99, 100,
		96, 0, 0, 0, 6, 0, 2, 48, 129, 208, 47, 0, 0, 0, 0, 73, 69,
		78, 68, 174, 66, 96, 130}
	r := bytes.NewReader(png1x1)
	http.ServeContent(c.Response(), c.Request(), "missing.png", time.Now(), r)
	return nil
}
