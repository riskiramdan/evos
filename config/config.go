package config

import (
	"fmt"
	"os"
)

const (
	dbhost     = "DB_HOST"
	dbdriver   = "DB_DRIVER"
	dbuser     = "DB_USER"
	dbpassword = "DB_PASSWORD"
	dbname     = "DB_NAME"
	dbport     = "DB_PORT"
)

// Config contains application configuration
type Config struct {
	DBConnectionString string
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	ImagePath          string
	AccountKey         string
	SecretKey          string
	CloudName          string
}

var config *Config

func getEnvOrDefault(env string, defaultVal string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultVal
	}
	return e
}

// GetConfiguration , get application configuration based on set environment
func GetConfiguration() (*Config, error) {
	if config != nil {
		return config, nil
	}

	dbDriver := getEnvOrDefault(dbdriver, "postgres")
	dbUser := getEnvOrDefault(dbuser, "postgres")
	dbPassword := getEnvOrDefault(dbpassword, "qweasd123")
	dbHost := getEnvOrDefault(dbhost, "127.0.0.1")
	dbPort := getEnvOrDefault(dbport, "5432")
	dbName := getEnvOrDefault(dbname, "evosdb")

	conStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", dbDriver, dbUser, dbPassword, dbHost, dbPort, dbName)
	// default configuration
	config := &Config{
		DBConnectionString: conStr,
	}
	return config, nil
}
