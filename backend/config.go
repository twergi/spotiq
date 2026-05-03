package main

import (
	"os"
	"strings"
)

type Config struct {
	BackendHostPort    string
	FrontURL           string
	RedirectURI        string
	ClientID           string
	ClientSecret       string
	CORSAllowedOrigins []string
}

func NewConfigFromEnv() *Config {
	conf := &Config{
		ClientID:           os.Getenv("CLIENT_ID"),
		ClientSecret:       os.Getenv("CLIENT_SECRET"),
		FrontURL:           os.Getenv("FRONTEND_URL"),
		RedirectURI:        os.Getenv("REDIRECT_URI"),
		CORSAllowedOrigins: strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		BackendHostPort:    os.Getenv("BACKEND_HOST_PORT"),
	}

	return conf
}
