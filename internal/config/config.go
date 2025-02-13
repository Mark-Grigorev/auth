package config

import (
	"log"
	"os"
	"strconv"

	"github.com/Mark-Grigorev/auth/internal/model"
)

const (
	defaultString = ""
	defaultInt    = 0
	defaultBool   = false
)

func Read() *model.Config {
	var config model.Config
	config.AppConfig = readAppConfig()
	config.DBConfig = readDBConfig()
	config.JWTConfig = readJWTConfig()
	config.RedisConfig = readRedisConfig()
	return &config
}

func readAppConfig() model.AppConfig {
	var config model.AppConfig
	config.Host = GetEnvAsType("HOST", defaultString)
	return config
}

func readDBConfig() model.DBConfig {
	var config model.DBConfig
	config.DBConnection = GetEnvAsType("DB_CONNECTION_STRING", defaultString)
	return config
}

func readJWTConfig() model.JWTConfig {
	var config model.JWTConfig
	config.SecretKey = GetEnvAsType("SECRET_KEY", defaultString)
	config.TokenDuration = GetEnvAsType("TOKEN_DURATION", defaultInt)
	return config
}

func readRedisConfig() model.RedisConfig {
	var config model.RedisConfig
	var defaultInt64 int64
	config.Servers = GetEnvAsType("REDIS_SERVERS", defaultString)
	config.Password = GetEnvAsType("REDIS_PASSWORD", defaultString)
	config.TTL = GetEnvAsType("REDIS_TOKEN_TTL", defaultInt64)
	return config
}

func GetEnvAsType[T any](key string, defaultValue T) T {
	var value string
	if value = os.Getenv(key); value == "" {
		log.Fatalf("не указан %s", key)
	}
	var result any
	switch any(defaultValue).(type) {
	case int:
		result, _ = strconv.Atoi(value)
	case bool:
		result, _ = strconv.ParseBool(value)
	case string:
		result = value
	case int64:
		result, _ = strconv.ParseInt(value, 10, 64)
	}

	return result.(T)
}
