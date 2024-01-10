package entities

import "os"

var RedirectURI string = os.Getenv("REDIRECT_URI")
var ClientID string = os.Getenv("CLIENT_ID")
var ClientSecret string = os.Getenv("CLIENT_SECRET")
