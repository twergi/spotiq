package main

type DeviceService struct {
	id    string
	name  string
	dtype string
	set   bool
}

func NewDeviceService() *DeviceService {
	return &DeviceService{}
}

func (d *DeviceService) Clear() {
	d.id = ""
	d.name = ""
	d.dtype = ""
	d.set = false
}

func (d *DeviceService) SetDevice(data Device) {
	d.id = data.ID
	d.name = data.Name
	d.dtype = data.Type
	d.set = true
}

func (d *DeviceService) Ready() bool {
	return d.set
}

func (d *DeviceService) GetCurrentDeviceID() string {
	return d.id
}
