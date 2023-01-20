package wellknown

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NodeInfo(c echo.Context) error {
	server := c.Get("server").(string)
	href := "https://" + c.Request().Host + "/nodeinfo/2.0"
	if server != "" {
		href += "?server=" + server
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"links": []map[string]string{
			{
				"rel":  "http://nodeinfo.diaspora.software/ns/schema/2.0",
				"href": href,
			},
		},
	})
}
