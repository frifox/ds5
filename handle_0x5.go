package ds5

import (
	"bytes"
	"encoding/binary"
)

const PS_FEATURE_CRC32_SEED = 0xA3

const DS_FEATURE_REPORT_CALIBRATION = 0x5
const DS_FEATURE_REPORT_CALIBRATION_SIZE = 41

type Input0x5 struct {
	ReportID byte

	InputGyro
	InputAccel

	Misc [2]byte

	CRC uint32
}
type InputGyro struct {
	PitchBias int16
	YawBias   int16
	RollBias  int16

	PitchPlus  int16
	PitchMinus int16

	YawPlus  int16
	YawMinus int16

	RollPlus  int16
	RollMinus int16

	SpeedPlus  int16
	SpeedMinus int16
}
type InputAccel struct {
	XPlus  int16
	XMinus int16

	YPlus  int16
	YMinus int16

	ZPlus  int16
	ZMinus int16
}

func (i *Input0x5) Unmarshal(data []byte) {
	//fmt.Printf("[FeatureReport0x5] Unmarshal len(%d) % X\n", len(data), data)

	reportReader := bytes.NewReader(data)
	err := binary.Read(reportReader, binary.LittleEndian, i)
	if err != nil {
		panic(err)
	}
}

// handle0x5 handles calibration info
func (d *Device) handle0x5(data []byte) {
	if ReportCRCIsValid(PS_FEATURE_CRC32_SEED, data) {
		d.Bus.Set("bt")
	} else {
		d.Bus.Set("usb")
	}

	r := Input0x5{}
	r.Unmarshal(data)

	d.Gyro.pitchCal.Bias = r.InputGyro.PitchBias
	d.Gyro.yawCal.Bias = r.InputGyro.YawBias
	d.Gyro.rollCal.Bias = r.InputGyro.RollBias

	d.Accel.cal = r.InputAccel

	//fmt.Printf("[%T] %#v\n", r, r)
}
