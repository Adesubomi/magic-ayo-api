package config

type DatabaseConfig struct {
	Connection string
	Host       string
	Port       string
	DbName     string `toml:"db_name"`
	User       string `toml:"user"`
	Password   string
}

type RedisConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string `toml:"user"`
	Password string
}

type SentryConfig struct {
	DNS     string `toml:"dns"`
	Debug   string `toml:"debug"`
	Release string `toml:"release"`
}

type LightningConfig struct {
	Url      string `toml:"url"`
	Port     string `toml:"port"`
	TlsCert  string `toml:"tls_cert"`
	Macaroon string `toml:"macaroon"`
}

type Config struct {
	AppName       string          `toml:"app_name"`
	AppPort       string          `toml:"app_port"`
	AppKey        string          `toml:"app_key"`
	Database      DatabaseConfig  `toml:"database"`
	Redis         RedisConfig     `toml:"redis"`
	Sentry        SentryConfig    `toml:"sentry"`
	LnSendNode    LightningConfig `toml:"lightning_send"`
	LnReceiveNode LightningConfig `toml:"lightning_receive"`
}
