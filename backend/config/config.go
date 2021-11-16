package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// TODO: Define config spec
}

func Init() *Config {
	godotenv.Load() // Handle err?
	c := &Config{}
	if err := envconfig.Process("APP", c); err != nil {
		log.Fatal(err) // TODO: use JSON log?
	}
	return c
}
