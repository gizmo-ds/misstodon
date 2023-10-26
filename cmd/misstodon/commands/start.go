package commands

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/gizmo-ds/misstodon/internal/api"
	"github.com/gizmo-ds/misstodon/internal/database"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/acme/autocert"
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
			Name: "fallback-server",
			Usage: "if proxy-server is not found in the request, the fallback server address will be used, " +
				`e.g. "misskey.io"`,
		},
	},
	Action: func(c *cli.Context) error {
		conf := global.Config
		global.DB = database.NewDatabase(
			conf.Database.Type,
			conf.Database.Address)
		defer global.DB.Close()
		if c.IsSet("fallbackServer") {
			conf.Proxy.FallbackServer = c.String("fallbackServer")
		}
		bindAddress, _ := utils.StrEvaluation(c.String("bind"), conf.Server.BindAddress)

		e := echo.New()
		e.HidePort, e.HideBanner = true, true

		api.Router(e)

		logStart := log.Info().Str("address", bindAddress)
		switch {
		case conf.Server.AutoTLS && conf.Server.Domain != "":
			e.Pre(middleware.HTTPSNonWWWRedirect())
			cacheDir, _ := filepath.Abs("./cert/.cache")
			e.AutoTLSManager.Cache = autocert.DirCache(cacheDir)
			e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(conf.Server.Domain)
			logStart.Msg("Starting server with AutoTLS")
			return e.StartAutoTLS(bindAddress)
		case conf.Server.TlsCertFile != "" && conf.Server.TlsKeyFile != "":
			logStart.Msg("Starting server with TLS")
			return e.StartTLS(bindAddress, conf.Server.TlsCertFile, conf.Server.TlsKeyFile)
		default:
			logStart.Msg("Starting server")
			return e.Start(bindAddress)
		}
	},
}
