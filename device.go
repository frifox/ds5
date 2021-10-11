package ds5

import (
	"fmt"
	"github.com/sstallion/go-hid"
)

const USB_VENDOR_ID_SONY = 0x54c
const USB_DEVICE_ID_SONY_PS5_CONTROLLER = 0xce6

type Device struct {
	Bus
	Touchpad Touchpad

	Buttons Buttons
	Axis    Axis
	Battery Battery

	LightBar   LightBar
	PlayerLEDs PlayerLEDs
	Mic        Mic
	Rumble     Rumble

	AliveFor
	OutputSequencer
	hid *hid.Device
}

func (d *Device) Find() (err error) {
	d.hid, err = hid.OpenFirst(USB_VENDOR_ID_SONY, USB_DEVICE_ID_SONY_PS5_CONTROLLER)

	if d.hid == nil {
		err = fmt.Errorf("%x:%dx not found: %v", USB_VENDOR_ID_SONY, USB_DEVICE_ID_SONY_PS5_CONTROLLER, err)
	}

	return
}

func (d *Device) Found() bool {
	return d.hid != nil
}

func (d *Device) Watch() {
	// by default DS5 sends 0x1 report, which is pretty basic.
	// requesting CALIBRATION report causes DS5 to send 0x31 report instead,
	// which includes goodies like Mute/Touch/Gyro/Accel/Battery/etc
	report0x5 := d.GetFeatureReport(DS_FEATURE_REPORT_CALIBRATION)
	switch len(report0x5) {
	case DS_FEATURE_REPORT_CALIBRATION_SIZE:
		d.handle0x5(report0x5)
	default:
		fmt.Printf("Unknown report0x5 len(%d)\n", len(report0x5))
		return
	}

	// pre-apply props
	d.ApplyProps(LightBarInit{})
	d.ApplyProps(d.LightBar, d.PlayerLEDs, d.Rumble)

	//d.GetFeatureReport(DS_FEATURE_REPORT_PAIRING_INFO, DS_FEATURE_REPORT_PAIRING_INFO_SIZE)
	//d.GetFeatureReport(DS_FEATURE_REPORT_FIRMWARE_INFO, DS_FEATURE_REPORT_FIRMWARE_INFO_SIZE)

	// read input reports 0x1/0x31 as they come in
	d.Reader()
}
func (d *Device) ApplyProps(props ...interface{}) {
	// validate props and save valid values to *Device for reference
	// then emit0x31() valid props for writing to DS5

	var applied []interface{}
	for _, prop := range props {
		switch p := prop.(type) {
		case LightBar:
			d.LightBar = p
			applied = append(applied, p)
		case PlayerLEDs:
			d.PlayerLEDs = p
			applied = append(applied, p)
		case Rumble:
			d.Rumble = p
			applied = append(applied, p)
		case Mic:
			d.Mic = p
			applied = append(applied, p)
		}
	}

	if len(applied) > 0 {
		switch d.Bus {
		case "usb":
			d.emit0x2(applied...)
		case "bt":
			d.emit0x31(applied...)
		}
	}
}
func (d *Device) GetFeatureReport(id uint8) []byte {
	report := make([]byte, 100)
	report[0] = id

	n, err := d.hid.GetFeatureReport(report)
	if err != nil {
		fmt.Printf("ERR GetFeatureReport(0x%X): %v\n", id, err)
	}

	// trim report what we actually read
	report = report[:n]

	return report
}

func (d *Device) Reader() {
	for {
		// BT report is the biggest input report we can expect
		report := make([]byte, DS_INPUT_REPORT_BT_SIZE)
		n, err := d.hid.Read(report)
		if err != nil {
			fmt.Printf("ERROR hid.Read(): %v\n", err)
			break
		}

		// trim report to bytes actually read
		report = report[:n]

		switch len(report) {
		case Input0x1ReportSize:
			d.handle0x1bt(report)
		case DS_INPUT_REPORT_USB_SIZE:
			d.handle0x1usb(report)
		case DS_INPUT_REPORT_BT_SIZE:
			d.handle0x31(report)
		default:
			fmt.Printf("[InputReport] UNKNOWN len(%d) % X\n", len(report), report)
		}

	}
}
