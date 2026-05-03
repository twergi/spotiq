package main

type QueueState struct {
	TokenSet  bool `json:"token_set"`
	DeviceSet bool `json:"device_set"`
}

type TokenPayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type AuthPayload struct {
	TokenPayload
	Error string `json:"error"`
}

type EnqueuePayload struct {
	URI string `json:"uri"`
}

type CurrentQueue struct {
	CurrentlyPlaying Song   `json:"currently_playing"`
	Queue            []Song `json:"queue"`
}

type Song struct {
	Name    string   `json:"name"`
	Album   Album    `json:"album"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	Name string `json:"name"`
}

type Album struct {
	Images []Image `json:"images"`
}

type Image struct {
	URL string `json:"url"`
}

type Device struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type DevicesPayload struct {
	Devices []Device `json:"devices"`
}

type SearchPayload struct {
	Tracks TracksPayload `json:"tracks"`
}

type TracksPayload struct {
	Items []Track `json:"items"`
}

type Track struct {
	Song
	EnqueuePayload
}
