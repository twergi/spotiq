package main

type Device struct {
	id    string
	name  string
	dtype string
}

func NewDevice() *Device {
	return &Device{}
}

type DeviceData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (d *Device) getData() (DeviceData, bool) {
	if !d.validateDevice() {
		return DeviceData{}, false
	}

	return DeviceData{
		ID:   d.id,
		Name: d.name,
		Type: d.dtype,
	}, true
}

func (d *Device) setData(data DeviceData) {
	d.id = data.ID
	d.name = data.Name
	d.dtype = data.Type
}

func (d *Device) validateDevice() bool {
	return d.id != ""
}
