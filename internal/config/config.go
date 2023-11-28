package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort                = "8080"
	defaultHTTPRWTimeout           = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes  = 1
	defaultDatabaseRefreshInterval = 30 * time.Second
)

type (
	Config struct {
		Postgres   PostgresConfig
		HTTP       HTTPConfig
		Cloudinary CloudinaryConfig
		Auth       AuthConfig
	}

	PostgresConfig struct {
		User     string
		Password string
		Port     string
		Host     string
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	}

	HTTPConfig struct {
		Host               string
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegaBytes int           `mapstructure:"maxHeaderMegaBytes"`
	}

	AuthConfig struct {
		JWT    JWTConfig
		Argon2 Argon2Config
	}

	JWTConfig struct {
		AccessTokenTTL time.Duration `mapstructure:"accessTokenTTL"`
		SigningKey     string
	}

	Argon2Config struct {
		MemoryMegaBytes uint32 `mapstructure:"memory"`
		Iterations      uint32 `mapstructure:"iterations"`
		Parallelism     uint8  `mapstructure:"parallelism"`
		SaltLength      uint32 `mapstructure:"saltLength"`
		KeyLength       uint32 `mapstructure:"keyLength"`
	}

	CloudinaryConfig struct {
		Cloud  string
		Key    string
		Secret string
	}
)

func InitConfig(configPath string) (*Config, error) {
	setDefaults()

	if err := parseConfigFile(configPath); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("hash", &cfg.Auth.Argon2); err != nil {
		return err
	}

	return viper.UnmarshalKey("http", &cfg.HTTP)
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")

	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	cfg.Cloudinary.Cloud = os.Getenv("CLOUDINARY_CLOUD")
	cfg.Cloudinary.Key = os.Getenv("CLOUDINARY_KEY")
	cfg.Cloudinary.Secret = os.Getenv("CLOUDINARY_SECRET")

	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}

func setDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.maxHeaderMegaBytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.readTimeout", defaultHTTPRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHTTPRWTimeout)
}
