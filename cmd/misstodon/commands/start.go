package commands

import (
	_ "embed"
	"fmt"

	"github.com/gizmo-ds/misstodon/internal/api"
	"github.com/gizmo-ds/misstodon/internal/database"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/internal/utils"
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
		appVersion := global.AppVersion
		if !c.Bool("no-color") {
			appVersion = "\033[1;31;40m" + appVersion + "\033[0m"
		}
		fmt.Printf("\n%s  %s\n\n", banner, appVersion)
		return nil
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "bind",
			Aliases: []string{"b"},
			Usage:   "bind address",
		},
		&cli.StringFlag{
			Name:  "url",
			Usage: "url of the server, used for generating links, " + `e.g. "https://example.com"`,
		},
		&cli.StringFlag{
			Name: "fallbackServer",
			Usage: "if proxy-server is not found in the request, the fallback server address will be used, " +
				`e.g. "misskey.io"`,
		},
	},
	Action: func(c *cli.Context) error {
		global.DB = database.NewDatabase(
			global.Config.Database.Type,
			global.Config.Database.Address)
		defer global.DB.Close()

		if c.IsSet("url") {
			global.Config.Server.Url = c.String("url")
		}
		if c.IsSet("fallbackServer") {
			global.Config.Proxy.FallbackServer = c.String("fallbackServer")
		}
		bindAddress, _ := utils.StrEvaluation(c.String("bind"), global.Config.Server.BindAddress)

		e := echo.New()
		e.HidePort, e.HideBanner = true, true
		api.Router(e)
		logStart := log.Info().Str("address", bindAddress)
		if global.Config.Server.TlsCertFile != "" && global.Config.Server.TlsKeyFile != "" {
			logStart.Msg("Starting server with TLS")
			return e.StartTLS(bindAddress,
				global.Config.Server.TlsCertFile, global.Config.Server.TlsKeyFile)
		}
		logStart.Msg("Starting server")
		return e.Start(bindAddress)
	},
}
