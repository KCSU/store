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
	OauthRedirectUrl string `split_words:"true"`
	LookupApiUrl     string `split_words:"true"`
	MailTemplateId   string `split_words:"true"`
	MailFrom         string `split_words:"true"`
	MailApiKey       string `split_words:"true"`
}

// Load configuration from environment variables or
// a .env file
func Init() *Config {
	_ = godotenv.Load() // TODO: Handle err (log no .env file)?
	c := &Config{}
	if err := envconfig.Process("", c); err != nil {
		log.Panic(err) // TODO: use JSON log?
	}
	return c
}
