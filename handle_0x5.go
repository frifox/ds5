package ds5

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const PS_FEATURE_CRC32_SEED = 0xA3

const DS_FEATURE_REPORT_CALIBRATION = 0x5
const DS_FEATURE_REPORT_CALIBRATION_SIZE = 41

const DS_ACC_RES_PER_G = 8192
const DS_ACC_RANGE = 4 * DS_ACC_RES_PER_G
const DS_GYRO_RES_PER_DEG_S = 1024
const DS_GYRO_RANGE = 2048 * DS_GYRO_RES_PER_DEG_S

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
	fmt.Printf("[FeatureReport0x5] Unmarshal len(%d) % X\n", len(data), data)

	reportReader := bytes.NewReader(data)
	err := binary.Read(reportReader, binary.LittleEndian, i)
	if err != nil {
		panic(err)
	}
}

func (d *Device) handle0x5(report []byte) {
	if ReportCRCIsValid(PS_FEATURE_CRC32_SEED, report) {
		d.Bus.Set("bt")
	} else {
		d.Bus.Set("usb")
	}

	r := Input0x5{}
	r.Unmarshal(report)

	return
	// TODO

	var pitchBias, yawBias, rollBias,
		pitchPlus, pitchMinus,
		yawPlus, yawMinus,
		rollPlus, rollMinus,
		speedPlus, speedMinus,
		accXPlus, accXMinus,
		accYPlus, accYMinus,
		accZPlus, accZMinus int16

	// r is 36 bytes long. Read 34 bytes into 17 int16 vars above
	// last 2 bytes are unknown...
	b := bytes.NewReader(report)

	binary.Read(b, binary.LittleEndian, &pitchBias)
	binary.Read(b, binary.LittleEndian, &yawBias)
	binary.Read(b, binary.LittleEndian, &rollBias)
	fmt.Printf("[Bias] Pitch:%+d Yaw:%+d Roll:%+d\n", pitchBias, yawBias, rollBias)

	binary.Read(b, binary.LittleEndian, &pitchPlus)
	binary.Read(b, binary.LittleEndian, &pitchMinus)
	fmt.Printf("[Pitch] %+d %+d\n", pitchPlus, pitchMinus)

	binary.Read(b, binary.LittleEndian, &yawPlus)
	binary.Read(b, binary.LittleEndian, &yawMinus)
	fmt.Printf("[Yaw] %+d %+d\n", yawPlus, yawMinus)

	binary.Read(b, binary.LittleEndian, &rollPlus)
	binary.Read(b, binary.LittleEndian, &rollMinus)
	fmt.Printf("[Roll] %+d %+d\n", rollPlus, rollMinus)

	binary.Read(b, binary.LittleEndian, &speedPlus)
	binary.Read(b, binary.LittleEndian, &speedMinus)
	fmt.Printf("[Speed] %+d %+d\n", speedPlus, speedMinus)

	// accel
	binary.Read(b, binary.LittleEndian, &accXPlus)
	binary.Read(b, binary.LittleEndian, &accXMinus)
	fmt.Printf("[AccelX] %+d %+d\n", accXPlus, accXMinus)

	binary.Read(b, binary.LittleEndian, &accYPlus)
	binary.Read(b, binary.LittleEndian, &accYMinus)
	fmt.Printf("[AccelY] %+d %+d\n", accYPlus, accYMinus)

	binary.Read(b, binary.LittleEndian, &accZPlus)
	binary.Read(b, binary.LittleEndian, &accZMinus)
	fmt.Printf("[AccelZ] %+d %+d\n", accZPlus, accZMinus)

	//var misc int16
	//binary.Read(b, binary.LittleEndian, &misc)
	//fmt.Printf("[Misc] %+d\n", misc)
}
