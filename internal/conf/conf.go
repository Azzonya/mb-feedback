package conf

import "github.com/caarlos0/env/v9"

var Conf = struct {
	HTTPListen string `env:"HTTP_LISTEN"`

	MbBrokerURL   string `env:"mb_broker_url"`
	MbBrokerToken string `env:"mb_broker_token"`

	VoximplantURL        string `env:"voximplant_url"`
	VoximplantToken      string `env:"voximplant_token"`
	VoximplantDomainName string `env:"voximplant_domain_name"`
	VoximplantTemplateID string `env:"voximplant_template_id"`
	VoximplantChannelID  string `env:"voximplant_channel_id"`

	PgDsn string `env:"pg_dsn"`
}{}

func init() {
	if err := env.Parse(&Conf); err != nil {
		panic(err)
	}
}
