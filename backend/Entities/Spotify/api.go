package Spotify

type SpotifyAPI interface {
	ClearData()
	GetData() APIData
	CreateTokenData(code string)
	GetDevices()
	SetDevice(deviceInfo DeviceData)
	GetCurrentQueue()
	AddToQueue(songInfo interface{})
	GetCurrentlyPlaying()
	Search(query map[string]string)
}

type SongData struct {
	URI string
}

type APIData struct {
	Token  TokenData
	Device DeviceData
}

type spotifyAPI struct {
	base_url string
	token    SpotifyToken
	device   SpotifyDevice
}

func (sAPI *spotifyAPI) ClearData() {
	sAPI.token = NewToken()
	sAPI.device = NewDevice()
}

func (sAPI *spotifyAPI) GetData() APIData {
	return APIData{
		Token:  sAPI.token.GetData(),
		Device: sAPI.device.GetData(),
	}
}

func (sAPI *spotifyAPI) CreateTokenData(code string) {

}

func (sAPI *spotifyAPI) GetDevices() {

}

func (sAPI *spotifyAPI) SetDevice(deviceData DeviceData) {

}

func (sAPI *spotifyAPI) GetCurrentQueue() {

}

func (sAPI *spotifyAPI) AddToQueue(songInfo interface{}) {

}

func (sAPI *spotifyAPI) GetCurrentlyPlaying() {

}

func (sAPI *spotifyAPI) Search(query map[string]string) {

}

func NewApi() SpotifyAPI {
	return &spotifyAPI{
		base_url: "https://api.spotify.com/v1",
		token:    NewToken(),
		device:   NewDevice(),
	}
}
