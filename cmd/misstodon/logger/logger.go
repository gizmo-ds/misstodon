package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gizmo-ds/misstodon/internal/global"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(noColor bool) {
	_ = os.MkdirAll(filepath.Dir(global.Config.Logger.Filename), 0750)
	zerolog.SetGlobalLevel(zerolog.Level(global.Config.Logger.Level))
	writers := []io.Writer{&lumberjack.Logger{
		Filename:   global.Config.Logger.Filename,
		MaxAge:     global.Config.Logger.MaxAge,
		MaxBackups: global.Config.Logger.MaxBackups,
	}}
	if global.Config.Logger.ConsoleWriter {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:           os.Stderr,
			NoColor:       noColor,
			TimeFormat:    "2006-01-02 15:04:05",
			FieldsExclude: []string{"stack"},
		})
	}
	log.Logger = zerolog.New(zerolog.MultiLevelWriter(writers...)).
		With().Timestamp().Stack().Logger()
}
