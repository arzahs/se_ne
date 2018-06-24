package core

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GenerateGoogleConfig(cfg *Config) *oauth2.Config{
	return &oauth2.Config{
		ClientID:     cfg.GoogleCid,
		ClientSecret: cfg.GoogleCsecret,
		RedirectURL:  "http://"+cfg.Email.Domain + "/api/v1/session/google/",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}