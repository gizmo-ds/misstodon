package global

import (
	"github.com/gizmo-ds/misstodon/internal/database"
	"github.com/jinzhu/configor"
)

type config struct {
	Proxy struct {
		FallbackServer string `toml:"fallback_server" yaml:"fallback_server"  env:"MISSTODON_FALLBACK_SERVER"`
	} `toml:"proxy" yaml:"proxy"`
	Server struct {
		BindAddress string `toml:"bind_address" yaml:"bind_address" env:"MISSTODON_SERVER_BIND_ADDRESS"`
		TlsCertFile string `toml:"tls_cert_file" yaml:"tls_cert_file" env:"MISSTODON_SERVER_TLS_CERT_FILE"`
		TlsKeyFile  string `toml:"tls_key_file" yaml:"tls_key_file" env:"MISSTODON_SERVER_TLS_KEY_FILE"`
	} `toml:"server" yaml:"server"`
	Logger struct {
		Level         int8   `toml:"level" yaml:"level" env:"MISSTODON_LOGGER_LEVEL"`
		ConsoleWriter bool   `toml:"console_writer" yaml:"console_writer" env:"MISSTODON_LOGGER_CONSOLE_WRITER"`
		RequestLogger bool   `toml:"request_logger" yaml:"request_logger" env:"MISSTODON_LOGGER_REQUEST_LOGGER"`
		Filename      string `toml:"filename" yaml:"filename" env:"MISSTODON_LOGGER_FILENAME"`
		MaxAge        int    `toml:"max_age" yaml:"max_age" env:"MISSTODON_LOGGER_MAX_AGE"`
		MaxBackups    int    `toml:"max_backups" yaml:"max_backups" env:"MISSTODON_LOGGER_MAX_BACKUPS"`
	} `toml:"logger" yaml:"logger"`
	Database struct {
		Type     database.DbType `toml:"type" yaml:"type" env:"MISSTODON_DATABASE_TYPE"`
		Address  string          `toml:"address" yaml:"address" env:"MISSTODON_DATABASE_ADDRESS"`
		Port     int             `toml:"port" yaml:"port" env:"MISSTODON_DATABASE_PORT"`
		User     string          `toml:"user" yaml:"user" env:"MISSTODON_DATABASE_USER"`
		Password string          `toml:"password" yaml:"password" env:"MISSTODON_DATABASE_PASSWORD"`
		DbName   string          `toml:"db_name" yaml:"db_name" env:"MISSTODON_DATABASE_DBNAME"`
	} `toml:"database" yaml:"database"`
}

var Config config

func LoadConfig(filename string) error {
	return configor.
		New(&configor.Config{Environment: "production"}).
		Load(&Config, filename)
}
