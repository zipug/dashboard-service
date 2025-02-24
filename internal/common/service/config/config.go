package config

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Redis struct {
	Host          string `toml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Port          int    `toml:"port" env:"REDIS_PORT" env-default:"6380"`
	DB            int    `toml:"db" env:"REDIS_DB" env-default:"0"`
	User          string `toml:"user" env:"REDIS_USER"`
	Password      string `toml:"password" env:"REDIS_USER_PASSWORD"`
	RedisPassword string `toml:"redis_password" env:"REDIS_PASSWORD"`
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
	Host              string        `toml:"host" env:"MINIO_HOST"`
	Port              int           `toml:"port" env:"MINIO_PORT"`
	User              string        `toml:"user" env:"MINIO_ROOT_USER"`
	Password          string        `toml:"password" env:"MINIO_ROOT_PASSWORD"`
	BucketArticles    string        `toml:"articles_bucket" env:"MINIO_ARTICLES_BUCKET"`
	BucketAttachments string        `toml:"attachments_bucket" env:"MINIO_ATTACHMENTS_BUCKET"`
	BucketAvatars     string        `toml:"avatars_bucket" env:"MINIO_AVATARS_BUCKET"`
	UrlLifetime       time.Duration `toml:"url_lifetime" env:"MINIO_URL_LIFETIME" env-default:"4h"`
	UseSsl            bool          `toml:"use_ssl" env:"MINIO_USE_SSL"`
}

type Server struct {
	Host               string        `toml:"host" env:"SERVER_HOST" env-default:"::1"`
	Port               int           `toml:"port" env:"SERVER_PORT" env-default:"4400"`
	DefaultApi         string        `toml:"default_api" env:"DEFAULT_API" env-default:"/api/v1/"`
	ReadTimeout        time.Duration `toml:"read_timeout" env:"READ_TIMEOUT" env-default:"5s"`
	WriteTimeout       time.Duration `toml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"5s"`
	IdleTimeout        time.Duration `toml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
	MaxParallelWorkers int           `toml:"max_parallel_workers" env:"MAX_PARALLEL_WORKERS" env-default:"3"`
	FrontEndUrl        string        `toml:"front_end_url" env:"FRONT_END_URL" env-default:"http://localhost:5173"`
}

type OTP struct {
	Max            int           `toml:"max" env:"OTP_MAX" env-default:"6"`
	ExpirationTime time.Duration `toml:"expiration_time" env:"OTP_EXPIRATION_TIME" env-default:"60s"`
}

type ENV string

type AppConfig struct {
	Env                   ENV           `env:"ENV" env-default:"production"`
	JwtSecretKey          string        `toml:"jwt_secret_key" env:"JWT_SECREY_KEY" env-required:"true"`
	AccessTokenExpiration time.Duration `toml:"access_token_exp" env:"ACCESS_TOKEN_EXP" env-required:"true"`
	Redis                 Redis         `toml:"redis"`
	Postgres              Postgres      `toml:"postgres"`
	Mongo                 Mongo         `toml:"mongo"`
	Prometheus            Prometheus    `toml:"prometheus"`
	MiniO                 MiniO         `toml:"minio"`
	Server                Server        `toml:"server" env-required:"true"`
	OTP                   OTP           `toml:"otp"`
	configPath            string
}

const (
	ENV_DEVELOPMENT ENV = "development"
	ENV_PRODUCTION  ENV = "production"
)

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
