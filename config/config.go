package config

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	MasterDB PostgresMasterConfig
	SlaveDB  PostgresSlaveConfig
	AWS      AwsConfig
	Logger   LoggerConfig
	Redis    RedisConfig
	Mail     MailConfig
	Admin    AdminConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	SSL          bool
	JWTSecretKey string
	Mode         string
}

type PostgresMasterConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	NameDB   string
}

type PostgresSlaveConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	NameDB   string
}

type LoggerConfig struct {
	Level       string
	Caller      bool
	Encoding    string
	Development bool
}

type AwsConfig struct {
	Endpoint       string
	MiniEndpoint   string
	MinioAccessKey string
	MinioSecretKey string
	UseSSL         bool
}

type RedisConfig struct {
	Host      string
	Port      int
	Password  string
	AuthDb    int
	SessionDb int
	LimiterDb int
}

type MailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type AdminConfig struct {
	User      string
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func NewAppConfig() (*Config, error) {

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	path := os.Getenv("CONFIG_PATH")

	if path != "" {
		// If CONFIG_PATH is a file → use it directly
		if filepath.Ext(path) != "" {
			v.SetConfigFile(path)
		} else {
			// If CONFIG_PATH is a directory → search inside it
			v.AddConfigPath(path)
		}
	} else {
		// fallback: project root
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
	}

	if err := v.ReadInConfig(); err != nil {
		var FileNotFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &FileNotFoundErr) {
			return nil, err
		}

		return nil, err
	}

	cfg := new(Config)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
