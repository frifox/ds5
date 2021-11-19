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
	MAC      MAC
	Version  Version
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

	// get MAC address, useful for automated BT pairing later on
	d.Reload0x9()

	// get HW/FW versions (available via BT only)
	d.Reload0x20()

	// by default DS5 sends 0x1 report, which is pretty basic.
	// requesting CALIBRATION report causes DS5 to send 0x31 report instead,
	// which includes goodies like Mute/Touch/Gyro/Accel/Battery/etc
	d.Reload0x5()

	// The hardware may have control over the LEDs (e.g. in Bluetooth on startup).
	// Reset the LEDs (lightbar, mute, player leds), so we can control them later.
	switch d.Bus.Type {
	case "usb":
		d.emit0x2(LEDSetup{})
	case "bt":
		d.emit0x31(LEDSetup{})
	}

	// will block until error
	d.Reader()

	// done. Close related workers
	d.Close()
}
func (d *Device) Reload0x5() (ok bool) {
	data := d.GetFeatureReport(DS_FEATURE_REPORT_CALIBRATION)
	switch len(data) {
	case DS_FEATURE_REPORT_CALIBRATION_SIZE:
		d.handle0x5(data)
		return true
	default:
		fmt.Printf("Unknown report0x5 len(%d)\n", len(data))
	}
	return false
}
func (d *Device) Reload0x9() (ok bool) {
	data := d.GetFeatureReport(DS_FEATURE_REPORT_PAIRING_INFO)
	switch len(data) {
	case DS_FEATURE_REPORT_PAIRING_INFO_SIZE:
		d.handle0x9(data)
		return true
	default:
		fmt.Printf("Unknown report0x9 len(%d)\n", len(data))
	}
	return false
}
func (d *Device) Reload0x20() (ok bool) {
	data := d.GetFeatureReport(DS_FEATURE_REPORT_FIRMWARE_INFO)
	switch len(data) {
	case DS_FEATURE_REPORT_FIRMWARE_INFO_SIZE:
		d.handle0x20(data)
		return true
	default:
		fmt.Printf("Unknown report0x20 len(%d)\n", len(data))
	}
	return false
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
		return nil
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
			//fmt.Printf("[DS.Write] %#v\n", report)

			data := report.Marshal()

			_, err := d.hid.Write(data)
			if err != nil {
				fmt.Printf("[%T] ERR hid.Write [%X]\n", report, data)
			}

		// ds5 KeepAlive
		case <-keepAlive.C:
			//fmt.Printf("KeepAlive: Reload0x5\n")
			if !d.Reload0x5() {
				//fmt.Printf("KeepAlive: Reload0x5 failed\n")
			}

		// shut down
		case <-d.Done():
			keepAlive.Stop()
			return
		}
	}
}
