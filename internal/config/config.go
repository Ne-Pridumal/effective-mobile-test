package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env            string `env:"ENV" evn-default:"local"`
	PassportApiUrl string `env:"PASSPORT_API_URL"`
	Postgres       `env:"POSTGRES" env-required:"true"`
	HttpServer     `env:"HTTP_SERVER"`
	Api            `env:"API"`
}

type HttpServer struct {
	Address     string        `env:"ADDRESS" env-default:"localhost"`
	Port        string        `env:"HTTP_PORT" env-default:"3000"`
	Timeout     time.Duration `env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

type Postgres struct {
	User     string `env:"USER" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	Address  string `env:"ADDRESS" env-required:"true"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	Ssl      string `env:"SSL" env-default:"disable"`
	Db       string `env:"DB" env-required:"true"`
}

type Api struct {
	Passport string `env:"PASSPORT"`
}

func MustLoad() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	configPath, isExist := os.LookupEnv("CONFIG_PATH")

	if configPath == "" && !isExist {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check does file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cnf Config
	if err := cleanenv.ReadConfig(configPath, &cnf); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cnf
}
