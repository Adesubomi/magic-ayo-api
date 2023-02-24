package config

type DatabaseConfig struct {
	Connection string
	Host       string
	Port       string
	DbName     string
	Username   string
	Password   string
}

type RedisConfig struct {
	Host     string
	Port     string
	DbName   string
	Username string
	Password string
}

type SentryConfig struct {
	DNS     string `env:"dns"`
	Debug   string `toml:"debug"`
	Release string `env:"release"`
}

type Config struct {
	AppName  string         `toml:"app_name"`
	AppPort  string         `toml:"app_port"`
	AppKey   string         `toml:"app_key"`
	Database DatabaseConfig `toml:"database"`
	Redis    RedisConfig    `toml:"redis"`
	Sentry   SentryConfig   `toml:"sentry"`
}
