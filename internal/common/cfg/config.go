package cfg

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env              string
	DB               DBConfig
	Server           ServerConfig
	Logger           LoggerConfig
	ElasticsearchURL string
	EmbedderURL      string
}

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port string
}

type LoggerConfig struct {
	Level      string
	OutputPath string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	cfg := &Config{
		Env: os.Getenv("ENV"),
		DB: DBConfig{
			Driver:   os.Getenv("DB_DRIVER"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		ElasticsearchURL: os.Getenv("ELASTICSEARCH_URL"),
		EmbedderURL:      os.Getenv("EMBEDDER_URL"),
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		Logger: LoggerConfig{
			Level:      os.Getenv("LOG_LEVEL"),
			OutputPath: os.Getenv("LOG_PATH"),
			MaxSize:    getEnvInt("LOG_MAX_SIZE", 100),
			MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 3),
			MaxAge:     getEnvInt("LOG_MAX_AGE", 28),
		},
	}

	return cfg, nil
}

func getEnvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}

	return i
}
