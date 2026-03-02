package app

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
	PG   PGConfig
	TG   Telegram
	LM   Limiter
}

type AppConfig struct {
	Name    string `env:"APP_NAME" env-default:"PolytechBot"`
	Version string `env:"APP_VERSION" env-default:"1.0.0"`
	Debug   bool   `env:"APP_DEBUG" env-default:"true"`
	BaseURL string `env:"BASE_URL"`
}

type HTTPConfig struct {
	Host       string `env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port       int    `env:"HTTP_PORT" env-default:"8080"`
	CORSOrigin string `env:"CORS_ORIGIN" env-default:"*"`
}

type PGConfig struct {
	DatabaseUser     string `env:"POSTGRES_USER"`
	DatabasePassword string `env:"POSTGRES_PASSWORD"`
	DatabaseHost     string `env:"POSTGRES_HOST"`
	DatabasePort     int    `env:"POSTGRES_PORT"`
	DatabaseDb       string `env:"POSTGRES_DB"`
}

type Telegram struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN"`
	WebUrl   string `env:"TELEGRAM_WEB_URL"`
}

type Limiter struct {
	RPS   float64 `env:"LIMITER_RPS" env-default:"10"`
	Burst int     `env:"LIMITER_BURST" env-default:"20"`
	Size  int     `env:"LIMITER_SIZE" env-default:"1000"`
}

func NewConfig() Config {
	cfg := Config{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			log.Fatal(err.Error())
		}
	}

	return cfg
}
