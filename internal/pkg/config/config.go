package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger  Logger     `yaml:"logger"`
	RServ   RESTServer `yaml:"server"`
	DB      DB         `yaml:"db"`
	RedisDB RedisDB    `yaml:"rdb"`
}

type RESTServer struct {
	Addr        string        `yaml:"addr"`
	ReadTimeout time.Duration `yaml:"readTimeout"`
	IdleTimeout time.Duration `yaml:"idleTimeout"`
}

type DB struct {
	Addr     string `yaml:"addr"`
	Username string `env:"POSTGRES_USER"     env-required:"true" yaml:"username"`
	Password string `env:"POSTGRES_PASSWORD" yaml:"password"`
	DB       string `env:"POSTGRES_DB"       env-required:"true" yaml:"db"`
	SSLmode  string `yaml:"sslmode"`
	MaxConns string `yaml:"maxConns"`
	Reload   bool   `yaml:"reload"`
	Version  int    `yaml:"version"`
}

type RedisDB struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Logger struct {
	Level     string   `yaml:"level"`
	Output    []string `yaml:"output"`
	ErrOutput []string `yaml:"errOutput"`
}

func New(configPath string) Config {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
