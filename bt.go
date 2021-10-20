package ds5

import (
	"github.com/muka/go-bluetooth/api"
	"strings"
)

func (d *Device) BTDisconnect() (err error) {
	var mac string

	if d.hid.Error() != nil {
		return
	}

	mac, err = d.hid.GetSerialNbr()
	if err != nil {
		return
	}

	adapter, err := api.GetDefaultAdapter()
	if err != nil {
		return err
	}

	devices, err := adapter.GetDevices()
	if err != nil {
		return nil
	}

	for _, device := range devices {
		props, err := device.GetProperties()
		if err != nil {
			continue
		}

		if !strings.EqualFold(props.Address, mac) {
			continue
		}

		err = device.Disconnect()
		if err != nil {
			return err
		}
	}

	return
}
