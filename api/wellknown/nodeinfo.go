package wellknown

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NodeInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"links": []map[string]string{
			{
				"rel": "http://nodeinfo.diaspora.software/ns/schema/2.0",
				"href": "https://" + c.Request().Host +
					"/nodeinfo/2.0?server=" +
					c.Get("server").(string),
			},
		},
	})
}
