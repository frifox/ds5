package ds5

type Version struct {
	HardwareVersion uint32
	FirmwareVersion uint32

	OnChange func(Version)
}

func (i *Version) Set(hw uint32, fw uint32) {
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
