package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"protravel-finance/pkg/logger"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Env string

var (
	Local Env = "local"
	Dev   Env = "dev"
	Prod  Env = "prod"
)

type Config struct {
	*jsonConfig
	*envConfig
}

type jsonConfig struct {
	Server  ServerConfig  `json:"server"`
	Handler HandlerConfig `json:"handler"`
}

type envConfig struct {
	Env      Env `env:"ENV" env-required:"true"`
	Postgres PostgresConfig
	Redis    RedisConfig

	Exchangerate ExchangerateConfig
}

type ServerConfig struct {
	Port           int           `json:"port"`
	ReadTimeout    time.Duration `json:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout"`
	MaxHeaderBytes int           `json:"max_header_bytes"`
}

type HandlerConfig struct {
	RequestTimeout  time.Duration `json:"request_timeout"`
	RegisterTimeout time.Duration `json:"register_timeout"`
}

type PostgresConfig struct {
	Host           string        `env:"POSTGRES_HOST" env-required:"true"`
	Port           int           `env:"POSTGRES_PORT" env-required:"true"`
	User           string        `env:"POSTGRES_USER" env-required:"true"`
	Password       string        `env:"POSTGRES_PASSWORD" env-required:"true"`
	Database       string        `env:"POSTGRES_DATABASE" env-required:"true"`
	IsDebug        bool          `env:"POSTGRES_IS_DEBUG" env-default:"false"`
	RequestTimeout time.Duration `env:"POSTGRES_REQUEST_TIMEOUT" env-required:"true"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-required:"true"`
	Port     string `env:"REDIS_PORT" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
	DB       int    `env:"REDIS_DB" env-required:"true"`
}

type ExchangerateConfig struct {
	BaseURL string `env:"EXCHANGERATE_BASE_URL" env-required:"true"`
}

func MustConfig(log logger.Logger) *Config {
	if err := godotenv.Load(); err != nil {
		log.Panic(fmt.Sprintf("failed download the file .env: %v", err))
	}

	path := fetchConfigPath()
	if path == "" {
		log.Panic("config path is empty")
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Panic("config file does not exist: " + path)
	}
	viper.AddConfigPath(filepath.Dir(path))
	viper.SetConfigType("json")
	viper.SetConfigName(filepath.Base(path))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	var jsonCfg jsonConfig

	err = viper.Unmarshal(&jsonCfg)
	if err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}
	validate := validator.New()

	err = validate.Struct(jsonCfg)
	if err != nil {
		log.Panicf("unable to validate config file: %v", err)
	}
	var envCfg envConfig

	err = cleanenv.ReadEnv(&envCfg)
	if err != nil {
		log.Panic("failed to read envConfig: " + err.Error())
	}
	return &Config{
		jsonConfig: &jsonCfg,
		envConfig:  &envCfg,
	}
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
