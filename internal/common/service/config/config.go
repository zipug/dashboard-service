package config

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Kafka struct {
	Host    string `toml:"host" env:"KAFKA_HOST"`
	Port    int    `toml:"port" env:"KAFKA_PORT"`
	Broker  string `toml:"broker" env:"KAFKA_BROKER"`
	Topic   string `toml:"topic" env:"KAFKA_TOPICS"`
	GroupID string `toml:"group_id" env:"KAFKA_GROUP_ID"`
}

type Redis struct {
	Host string `toml:"host" env:"REDIS_HOST"`
	Port int    `toml:"port" env:"REDIS_PORT"`
}

type Postgres struct {
	Host           string `toml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port           int    `toml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	User           string `toml:"user" env:"POSTGRES_USER"`
	Password       string `toml:"password" env:"POSTGRES_PASSWORD"`
	DBName         string `toml:"db_name" env:"POSTGRES_DB_NAME"`
	SSLMode        string `toml:"ssl_mode" env:"POSTGRES_SSL_MODE" env-default:"disable"`
	MigrationsPath string `toml:"migrations_path" env:"POSTGRES_MIGRATIONS_PATH" env-required:"true"`
}

type Mongo struct {
	Host string `toml:"host" env:"MONGO_HOST"`
	Port int    `toml:"host" env:"MONGO_PORT"`
}

type Prometheus struct {
	Host string `toml:"host" env:"PROMETHEUS_HOST"`
	Port int    `toml:"host" env:"PROMETHEUS_PORT"`
}

type MiniO struct {
	Host   string `toml:"host" env:"MINIO_HOST"`
	Port   int    `toml:"port" env:"MINIO_PORT"`
	Bucket string `toml:"bucket" env:"MINIO_BUCKET"`
}

type Server struct {
	Host         string        `toml:"host" env:"SERVER_HOST" env-default:"::1"`
	Port         int           `toml:"port" env:"SERVER_PORT" env-default:"4400"`
	DefaultApi   string        `toml:"default_api" env:"DEFAULT_API" env-default:"/api/v1/"`
	ReadTimeout  time.Duration `toml:"read_timeout" env:"READ_TIMEOUT" env-default:"5s"`
	WriteTimeout time.Duration `toml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"5s"`
	IdleTimeout  time.Duration `toml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
}

type OTP struct {
	Max int `toml:"max" env:"OTP_MAX" env-default:"6"`
}

type AppConfig struct {
	Env                   string        `env:"ENV" env-default:"development"`
	JwtSecretKey          string        `toml:"jwt_secret_key" env:"JWT_SECREY_KEY" env-required:"true"`
	AccessTokenExpiration time.Duration `toml:"access_token_exp" env:"ACCESS_TOKEN_EXP" env-required:"true"`
	Kafka                 Kafka         `toml:"kafka"`
	Redis                 Redis         `toml:"redis"`
	Postgres              Postgres      `toml:"postgres"`
	Mongo                 Mongo         `toml:"mongo"`
	Prometheus            Prometheus    `toml:"prometheus"`
	MiniO                 MiniO         `toml:"minio"`
	Server                Server        `toml:"server" env-required:"true"`
	OTP                   OTP           `toml:"otp"`
	configPath            string
}

func NewConfigService() *AppConfig {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		if cfg_path, ok := os.LookupEnv("CONFIG_PATH"); ok {
			path = cfg_path
		}
	}

	cfg := &AppConfig{configPath: path}
	cfg.load()

	return cfg
}

func (cfg *AppConfig) load() error {
	if cfg.configPath == "" {
		return errors.New("config path is not set")
	}

	if err := cleanenv.ReadConfig(cfg.configPath, cfg); err != nil {
		return err
	}

	return nil
}
