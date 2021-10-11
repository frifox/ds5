package ds5

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

const PS_INPUT_CRC32_SEED = 0xA1

const DS_INPUT_REPORT_BT = 0x31
const DS_INPUT_REPORT_BT_SIZE = 78

type Input0x31 struct {
	ReportID byte
	Unknown  byte

	LeftX, LeftY   byte
	RightX, RightY byte
	L2, R2         byte

	SeqNumber byte

	Buttons  [4]byte
	Reserved [4]byte

	Gyro      [3]uint16
	Accel     [3]uint16
	Timestamp uint32
	Reserved2 byte

	TouchPoints [2]uint32

	Reserved3 [12]byte
	Status    byte
	Reserved4 [10]byte

	Reserved5 [9]byte
	CRC       uint32
}

func (r *Input0x31) Unmarshal(data []byte) {
	reportReader := bytes.NewReader(data)
	err := binary.Read(reportReader, binary.LittleEndian, r)
	if err != nil {
		panic(err)
	}
}
func (r *Input0x31) Marshal() []byte {
	r.CRC = ReportCRC(PS_INPUT_CRC32_SEED, r)

	buff := bytes.Buffer{}
	if err := binary.Write(&buff, binary.LittleEndian, r); err != nil {
		panic(err)
	}

	// integrity check
	if buff.Len() != DS_INPUT_REPORT_BT_SIZE {
		panic(fmt.Sprintf("[%T] len(%d) != %d", r, buff.Len(), DS_INPUT_REPORT_BT_SIZE))
	}

	return buff.Bytes()
}

func (d *Device) handle0x31(report []byte) {
	if !ReportCRCIsValid(PS_INPUT_CRC32_SEED, report) {
		fmt.Printf("ERR CRC check failed [% X]\n", report)
		return
	}

	r := Input0x31{}
	r.Unmarshal(report)

	// sanity check
	if !bytes.Equal(report, r.Marshal()) {
		fmt.Printf("ERR Input0x31 sanity check [%X] != [%X], %+v\n", report, r.Marshal(), r)
	}

	d.Axis.Left.Set(r.LeftX, r.LeftY)
	d.Axis.Right.Set(r.RightX, r.RightY)

	d.Axis.L2.Set(r.L2)
	d.Axis.R2.Set(r.R2)

	switch r.Buttons[0] & 0xf {
	case 0x0:
		d.Buttons.DPadUp.Set(true)
	case 0x1:
		d.Buttons.DPadUp.Set(true)
		d.Buttons.DPadRight.Set(true)
	case 0x2:
		d.Buttons.DPadRight.Set(true)
	case 0x3:
		d.Buttons.DPadRight.Set(true)
		d.Buttons.DPadDown.Set(true)
	case 0x4:
		d.Buttons.DPadDown.Set(true)
	case 0x5:
		d.Buttons.DPadDown.Set(true)
		d.Buttons.DPadLeft.Set(true)
	case 0x6:
		d.Buttons.DPadLeft.Set(true)
	case 0x7:
		d.Buttons.DPadLeft.Set(true)
		d.Buttons.DPadUp.Set(true)
	case 0x8:
		d.Buttons.DPadUp.Set(false)
		d.Buttons.DPadRight.Set(false)
		d.Buttons.DPadDown.Set(false)
		d.Buttons.DPadLeft.Set(false)
	}
	d.Buttons.Square.Set(r.Buttons[0]>>4&1 == 1)
	d.Buttons.Cross.Set(r.Buttons[0]>>5&1 == 1)
	d.Buttons.Circle.Set(r.Buttons[0]>>6&1 == 1)
	d.Buttons.Triangle.Set(r.Buttons[0]>>7&1 == 1)

	d.Buttons.L1.Set(r.Buttons[1]>>0&1 == 1)
	d.Buttons.R1.Set(r.Buttons[1]>>1&1 == 1)
	d.Buttons.L2.Set(r.Buttons[1]>>2&1 == 1)
	d.Buttons.R2.Set(r.Buttons[1]>>3&1 == 1)
	d.Buttons.R2.Set(r.Buttons[1]>>3&1 == 1)
	d.Buttons.Share.Set(r.Buttons[1]>>4&1 == 1)
	d.Buttons.Options.Set(r.Buttons[1]>>5&1 == 1)
	d.Buttons.Left.Set(r.Buttons[1]>>6&1 == 1)
	d.Buttons.Right.Set(r.Buttons[1]>>7&1 == 1)

	d.Buttons.PS.Set(r.Buttons[2]>>0&1 == 1)
	d.Buttons.Touchpad.Set(r.Buttons[2]>>1&1 == 1)
	d.Buttons.Mute.Set(r.Buttons[2]>>2&1 == 1)

	// TODO Gyro, Accel

	// 1 r.Timestamp = 1/3 microseconds
	aliveFor := time.Duration(float64(time.Microsecond) * float64(r.Timestamp) / 3)
	aliveFor = aliveFor.Round(time.Second)
	d.AliveFor.Set(aliveFor)

	{
		touch := r.TouchPoints[0] // 32-bit value (4 bytes)

		info := uint8(touch & 0xff) // first byte
		xy := touch >> 8            // next 3 bytes

		id := info >> 0 & 0x7f   // 7 bits long
		active := info>>7&1 == 0 // 0 = active, 1 = inactive
		x := xy & 0x000fff       // 12-bit value on right
		y := xy & 0xfff000 >> 12 // 12-bit value on left

		d.Touchpad.Touch1.Set(id, active, int(x), int(y))
	}

	{
		touch := r.TouchPoints[1] // 32 bit value (4 bytes)

		info := uint8(touch & 0xff) // first byte
		xy := touch >> 8            // next 3 bytes

		id := info >> 0 & 0x7f   // 7 bits long
		active := info>>7&1 == 0 // 0 = active, 1 = inactive
		x := xy & 0x000fff       // 12-bit value on right
		y := xy & 0xfff000 >> 12 // 12-bit value on left

		d.Touchpad.Touch2.Set(id, active, int(x), int(y))
	}

	// Each unit of battery data corresponds to 10%
	// 0 = 0-9%, 1 = 10-19%, .. 10 = 100%
	percent := (r.Status & 0xf) * 10
	switch r.Status >> 4 {
	case 0x0:
		d.Battery.Set("Discharging", percent)
	case 0x1:
		d.Battery.Set("Charging", percent)
	case 0x2:
		d.Battery.Set("Full", 100)
	case 0xa:
		d.Battery.Set("Volt/Temp OutOfRange", 0)
	case 0xb:
		d.Battery.Set("Temperature ERR", 0)
	case 0xf:
		d.Battery.Set("Charging ERR", 0)
	default:
		d.Battery.Set("Unknown", 0)
	}
}
