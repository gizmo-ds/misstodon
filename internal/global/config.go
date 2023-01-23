package global

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gizmo-ds/misstodon/internal/database"
)

type config struct {
	Proxy struct {
		FallbackServer string `toml:"fallback_server"`
	} `toml:"proxy"`
	Server struct {
		Url         string `toml:"url"`
		BindAddress string `toml:"bind_address"`
		TlsCertFile string `toml:"tls_cert_file"`
		TlsKeyFile  string `toml:"tls_key_file"`
	} `toml:"server"`
	Logger struct {
		Level         int8   `toml:"level"`
		ConsoleWriter bool   `toml:"console_writer"`
		RequestLogger bool   `toml:"request_logger"`
		Filename      string `toml:"filename"`
		MaxAge        int    `toml:"max_age"`
		MaxBackups    int    `toml:"max_backups"`
	} `toml:"logger"`
	Database struct {
		Type     database.DbType `toml:"type"`
		Address  string          `toml:"address"`
		Port     int             `toml:"port"`
		User     string          `toml:"user"`
		Password string          `toml:"password"`
		Database string          `toml:"database"`
	} `toml:"database"`
}

var Config config

func LoadConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, &Config)
}
