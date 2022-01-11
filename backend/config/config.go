package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Configuration data for the app
//
// TODO: split into separate configs?
type Config struct {
	Debug            bool
	DbConnection     string `split_words:"true"`
	CookieSecret     string `split_words:"true"`
	JwtSecret        string `split_words:"true"`
	OauthClientKey   string `split_words:"true"`
	OauthSecretKey   string `split_words:"true"`
	OauthCallbackUrl string `split_words:"true"`
	LookupApiUrl     string `split_words:"true"`
}

// Load configuration from environment variables or
// a .env file
func Init() *Config {
	godotenv.Load() // Handle err?
	c := &Config{}
	if err := envconfig.Process("APP", c); err != nil {
		log.Panic(err) // TODO: use JSON log?
	}
	return c
}
