package misskey_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	testServer string
	testToken  string
	testAcct   string
)

func TestMain(m *testing.M) {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
	})
	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("failed to load .env file")
		return
	}
	testServer = os.Getenv("TEST_SERVER")
	testToken = os.Getenv("TEST_TOKEN")
	testAcct = os.Getenv("TEST_ACCT")
	m.Run()
}
