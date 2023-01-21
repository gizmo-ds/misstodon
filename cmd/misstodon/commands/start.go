package commands

import (
	_ "embed"
	"fmt"

	"github.com/gizmo-ds/misstodon/api"
	"github.com/gizmo-ds/misstodon/internal/database"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

//go:embed banner.txt
var banner string

var Start = &cli.Command{
	Name:  "start",
	Usage: "Start the server",
	Before: func(c *cli.Context) error {
		fmt.Printf("\n%s  \033[1;31;40m%s\033[0m\n\n", banner, global.AppVersion)
		return nil
	},
	Action: func(c *cli.Context) error {
		global.DB = database.NewDatabase(
			global.Config.Database.Type,
			global.Config.Database.Address)
		defer global.DB.Close()

		e := echo.New()
		e.HidePort, e.HideBanner = true, true
		api.Router(e)
		l := log.Info().Str("address", global.Config.Server.BindAddress)
		if global.Config.Server.TlsCertFile != "" && global.Config.Server.TlsKeyFile != "" {
			l.Msg("Starting server with TLS")
			return e.StartTLS(global.Config.Server.BindAddress,
				global.Config.Server.TlsCertFile, global.Config.Server.TlsKeyFile)
		}
		l.Msg("Starting server")
		return e.Start(global.Config.Server.BindAddress)
	},
}
