package ds5

const DS_FEATURE_REPORT_FIRMWARE_INFO = 0x20
const DS_FEATURE_REPORT_FIRMWARE_INFO_SIZE = 64

func (d *Device) handle0x20(report []byte) {
	// hid-playstation.c: dualsense_get_firmware_info()

	//ds->base.hw_version = get_unaligned_le32(&buf[24]);
	//ds->base.fw_version = get_unaligned_le32(&buf[28]);
}
