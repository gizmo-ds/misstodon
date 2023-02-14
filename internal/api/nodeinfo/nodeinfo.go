package nodeinfo

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/api/middleware"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Group) {
	group := e.Group("/nodeinfo", middleware.CORS())
	group.GET("/2.0", InfoHandler)
}

func InfoHandler(c echo.Context) error {
	server := c.Get("server").(string)
	var err error
	info := models.NodeInfo{
		Version: "2.0",
		Software: models.NodeInfoSoftware{
			Name:    "misstodon",
			Version: global.AppVersion,
		},
		Protocols: []string{"activitypub"},
		Services: models.NodeInfoServices{
			Inbound:  []string{},
			Outbound: []string{},
		},
		Metadata: struct{}{},
	}
	if server != "" {
		info, err = misskey.NodeInfo(
			server,
			models.NodeInfo{
				Version: "2.0",
				Software: models.NodeInfoSoftware{
					Name:    "misstodon",
					Version: global.AppVersion,
				},
				Protocols: []string{"activitypub"},
				Services: models.NodeInfoServices{
					Inbound:  []string{},
					Outbound: []string{},
				},
			})
		if err != nil {
			return err
		}
	}
	return c.JSON(http.StatusOK, info)
}
