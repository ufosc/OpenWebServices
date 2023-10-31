package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config encapsulates environment variables.
type Config struct {
	GIN_MODE         string
	MONGO_URI        string
	DB_NAME          string
	NOTIF_EMAIL_ADDR string
	FE_PROXY_PORT    string
}

// GetDefaultConfig populates a Config instance with default configuration
// options.
func GetDefaultConfig() (c Config) {
	c.GIN_MODE = "debug"
	c.MONGO_URI = ""
	c.DB_NAME = "testing"
	c.NOTIF_EMAIL_ADDR = "no-reply.notifications@ufosc.org"
	c.FE_PROXY_PORT = ""
	return c
}

// GetConfig populates a Config instance with default configuration options
// or options specified in OS variables or the .env file, when applicable.
func GetConfig() Config {
	c := GetDefaultConfig()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if mode := os.Getenv("GIN_MODE"); mode == "debug" || mode == "production" {
		c.GIN_MODE = mode
	}
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		c.MONGO_URI = uri
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		c.DB_NAME = dbname
	}
	if addr := os.Getenv("NOTIF_EMAIL_ADDR"); addr != "" {
		c.NOTIF_EMAIL_ADDR = addr
	}
	if proxy := os.Getenv("FE_PROXY_PORT"); proxy != "" {
		c.FE_PROXY_PORT = proxy
	}

	return c
}
