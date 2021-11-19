package ds5

type Info struct {
	MAC string

	HardwareVersion uint32
	FirmwareVersion uint32

	OnChange func(Info)
}

func (i *Info) SetMAC(mac string) {
	if i.MAC == mac {
		return // nothing changed
	}

	i.MAC = mac

	// any callbacks?
	if i.OnChange != nil {
		go i.OnChange(*i)
	}
}

func (i *Info) SetHW(hw uint32, fw uint32) {
	if i.HardwareVersion == hw && i.FirmwareVersion == fw {
		return // nothing changed
	}

	i.HardwareVersion = hw
	i.FirmwareVersion = fw

	// any callbacks?
	if i.OnChange != nil {
		go i.OnChange(*i)
	}
}
