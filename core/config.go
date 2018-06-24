package core

import (
	"github.com/joeshaw/envdecode"
	"se_ne/db"
)

type Config struct {
	DB            db.Config
	Email         EmailConfig
	Host          string `env:"APP_HOST,default=localhost"`
	Port          int    `env:"APP_PORT,default=8000"`
	StaticPath    string `env:"APP_STATIC_PATH,default=./assets/"`
	TemplatePath  string `env:"APP_TEMPLATE_PATH,default=./templates/"`
	SecretKey     string `env:"APP_SECRET_KEY,default=my-secret-key"`
	GoogleAPIKey  string `env:"APP_GOOGLE_MAP_API_KEY,default=AIzaSyBkQbmkLbtC4HgaM41BLx8Qqb-djFjTX5E"`
	GoogleCid     string `env:"APP_GOOGLE_CID,required"`
	GoogleCsecret string `env:"APP_GOOGLE_CSECRET,required"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{DB: db.Config{}}
	err := envdecode.Decode(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
