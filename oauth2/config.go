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
	SECRET           string
	PORT             string
	WEBSMTP          string
}

// GetDefaultConfig populates a Config instance with default configuration
// options.
func GetDefaultConfig() (c Config) {
	c.GIN_MODE = "debug"
	c.MONGO_URI = ""
	c.DB_NAME = "testing"
	c.NOTIF_EMAIL_ADDR = "no-reply.notifications@ufosc.org"
	c.SECRET = "369369369369369369"
	c.PORT = "8080"
	c.WEBSMTP = "http://localhost:3001"
	return c
}

// GetConfig populates a Config instance with default configuration options
// or options specified in OS variables or the .env file, when applicable.
func GetConfig() Config {
	c := GetDefaultConfig()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if mode := os.Getenv("GIN_MODE"); mode == "debug" || mode == "release" {
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
	if secret := os.Getenv("SECRET"); secret != "" {
		c.SECRET = secret
	}
	if port := os.Getenv("PORT"); port != "" {
		c.PORT = port
	}
	if websmtp := os.Getenv("WEBSMTP"); websmtp != "" {
		c.WEBSMTP = websmtp
	}

	return c
}
