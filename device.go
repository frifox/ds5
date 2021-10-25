package ds5

import (
	"context"
	"fmt"
	"github.com/sstallion/go-hid"
	"time"
)

const USB_VENDOR_ID_SONY = 0x54c
const USB_DEVICE_ID_SONY_PS5_CONTROLLER = 0xce6

type Device struct {
	Bus      Bus
	AliveFor AliveFor
	Battery  Battery

	Buttons  Buttons
	Axis     Axis
	Touchpad Touchpad
	Gyro     Gyro
	Accel    Accel

	LightBar   LightBar
	PlayerLEDs PlayerLEDs
	Mic        Mic
	Rumble     Rumble

	outputSequencer outputSequencer
	hid             *hid.Device

	writer chan Report

	context.Context
	Close context.CancelFunc
}

type Report interface {
	Marshal() []byte
}

type LEDSetup struct{}

func (d *Device) Find() (err error) {
	d.hid, err = hid.OpenFirst(USB_VENDOR_ID_SONY, USB_DEVICE_ID_SONY_PS5_CONTROLLER)

	if d.hid == nil {
		err = fmt.Errorf("%x:%dx not found: %v", USB_VENDOR_ID_SONY, USB_DEVICE_ID_SONY_PS5_CONTROLLER, err)
	}

	d.Context, d.Close = context.WithCancel(context.Background())

	return
}

func (d *Device) Found() bool {
	return d.hid != nil
}

func (d *Device) Run() {
	d.writer = make(chan Report)
	go d.Writer()

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

	// The hardware may have control over the LEDs (e.g. in Bluetooth on startup).
	// Reset the LEDs (lightbar, mute, player leds), so we can control them from software.
	switch d.Bus.Type {
	case "usb":
		d.emit0x2(LEDSetup{})
	case "bt":
		d.emit0x31(LEDSetup{})
	}

	//d.GetFeatureReport(DS_FEATURE_REPORT_PAIRING_INFO, DS_FEATURE_REPORT_PAIRING_INFO_SIZE)
	//d.GetFeatureReport(DS_FEATURE_REPORT_FIRMWARE_INFO, DS_FEATURE_REPORT_FIRMWARE_INFO_SIZE)

	// will block until error
	d.Reader()

	// done. Close related workers
	d.Close()
}
func (d *Device) ApplyProps() {
	switch d.Bus.Type {
	case "usb":
		d.emit0x2()
	case "bt":
		d.emit0x31()
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

func (d *Device) Writer() {
	keepAlive := time.NewTicker(time.Minute)

	for {
		select {
		case report := <-d.writer:
			data := report.Marshal()

			_, err := d.hid.Write(data)
			if err != nil {
				fmt.Printf("[%T] ERR hid.Write | %v |len(%d) [%X]\n", report, err, data, data)
			} else {
				//fmt.Printf("[Emit0x31 #%d] Send %d Bytes. Len(%d) [%X]\n", goID(), n, len(data), data)
			}

		// ds5 KeepAlive
		case <-keepAlive.C:
			fmt.Printf("KeepAlive\n")
			d.ApplyProps()

		// shut down
		case <-d.Done():
			keepAlive.Stop()
			return
		}
	}
}
