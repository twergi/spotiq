package Spotify

type DeviceData struct {
	ID         string
	Name       string
	DeviceType string
}

type SpotifyDevice interface {
	GetData() DeviceData
	SetData(deviceData DeviceData)
	IsValid() bool
}

type spotifyDevice struct {
	data DeviceData
}

func (sDevice *spotifyDevice) GetData() DeviceData {
	return sDevice.data
}
func (sDevice *spotifyDevice) SetData(deviceData DeviceData) {
	sDevice.data = deviceData
}
func (sDevice *spotifyDevice) IsValid() bool {
	return sDevice.data.ID != ""
}

func NewDevice() SpotifyDevice {
	return &spotifyDevice{}
}
