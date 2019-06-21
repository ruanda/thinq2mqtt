package thinq

type DeviceService struct {
	service
	devices []Device
}

type Device struct {
	client *Client
	Data   *DeviceData
}

func (s *DeviceService) List() []Device {
	return s.devices
}

func (s *DeviceService) Add(deviceData *DeviceData) {
	for _, d := range s.devices {
		if d.Data.ID == deviceData.ID {
			return
		}
	}
	// TODO: add locking
	s.devices = append(s.devices, Device{
		client: s.client,
		Data:   deviceData,
	})
}
