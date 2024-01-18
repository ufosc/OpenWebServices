package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config encapsulates environment variables.
type Config struct {
	GIN_MODE string // "release" or "debug".
	PORT     string // server port.
	THREADS  string // number of worker threads.
}

// GetDefaultConfig populates a Config strut with default
// configuration options.
func GetDefaultConfig() (c Config) {
	c.GIN_MODE = "debug"
	c.PORT = "8080"
	c.THREADS = "1"
	return c
}

// GetConfig populates a Config struct with default configuration
// options, options specified in OS variables, or options specified
// in the .env file (when available).
func GetConfig() Config {
	c := GetDefaultConfig()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if mode := os.Getenv("GIN_MODE"); mode == "debug" || mode == "release" {
		c.GIN_MODE = mode
	}
	if port := os.Getenv("PORT"); port != "" {
		c.PORT = port
	}
	if threads := os.Getenv("THREADS"); threads != "" {
		c.THREADS = threads
	}

	return c
}
