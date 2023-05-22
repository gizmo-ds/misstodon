package main

import (
	_ "embed"
	"os"

	"github.com/gizmo-ds/misstodon/cmd/misstodon/commands"
	"github.com/gizmo-ds/misstodon/cmd/misstodon/logger"
	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/gizmo-ds/misstodon/proxy/misskey"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/urfave/cli/v2"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	err := (&cli.App{
		Name:                 "misstodon",
		Usage:                "misskey api proxy",
		Version:              global.AppVersion,
		EnableBashCompletion: true,
		Suggest:              true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "Disable color output",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "config file",
				Value:   "config.toml",
			},
		},
		Before: func(c *cli.Context) error {
			if err := global.LoadConfig(c.String("config")); err != nil {
				log.Fatal().Stack().Err(errors.WithStack(err)).Msg("Failed to load config")
			}
			logger.Init(c.Bool("no-color"))
			misskey.SetHeader("User-Agent", "misstodon/"+global.AppVersion)
			return nil
		},
		Commands: []*cli.Command{
			commands.Start,
		},
	}).Run(os.Args)
	if err != nil {
		log.Fatal().Err(errors.WithStack(err)).Msg("Failed to start")
	}
}
