package global

import (
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	ServerBind string `toml:"server_bind"`
	Server     struct {
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
}

var Config config

func LoadConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, &Config)
}
