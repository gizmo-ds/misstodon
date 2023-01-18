package nodeinfo

import (
	"net/http"

	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/models"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/labstack/echo/v4"
)

func NodeInfo(c echo.Context) error {
	info, err := misskey.NodeInfo(
		c.Get("server").(string),
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
	return c.JSON(http.StatusOK, info)
}
