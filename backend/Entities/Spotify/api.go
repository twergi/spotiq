package Spotify

type SpotifyAPI interface {
	ClearData()
	GetData()
	CreateTokenData()
	GetDevices()
	SetDevice(deviceInfo interface{})
	GetCurrentQueue()
	AddToQueue(songInfo interface{})
	GetCurrentlyPlaying()
	Search(query map[string]string)
}

type spotifyAPI struct {
	base_url string
	token    SpotifyToken
	decice   SpotifyDevice
}
