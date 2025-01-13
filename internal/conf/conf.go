package conf

import "github.com/caarlos0/env/v9"

// Conf represents the application configuration.
var Conf = struct {
	ServerAddress string `env:"server_address"`

	PgDsn string `env:"DATABASE_URI"`
}{}

// init initializes the configuration by parsing environment variables.
// It sets default values for fields.
func init() {
	if err := env.Parse(&Conf); err != nil {
		panic(err)
	}
}
