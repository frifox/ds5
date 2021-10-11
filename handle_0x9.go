package ds5

const DS_FEATURE_REPORT_PAIRING_INFO = 0x09
const DS_FEATURE_REPORT_PAIRING_INFO_SIZE = 20

func (d *Device) handle0x9(report []byte) {
	// hid-playstation.c: dualsense_get_mac_address()

	//memcpy(ds->base.mac_address, &buf[1], sizeof(ds->base.mac_address));
}
