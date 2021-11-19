package ds5

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const DS_FEATURE_REPORT_FIRMWARE_INFO = 0x20
const DS_FEATURE_REPORT_FIRMWARE_INFO_SIZE = 64

type Input0x20 struct {
	ReportID byte

	Unknown1 [23]byte

	HardwareVersion uint32
	FirmwareVersion uint32

	Unknown2 [28]byte
	CRC      uint32
}

func (r *Input0x20) Unmarshal(data []byte) {
	reportReader := bytes.NewReader(data)
	err := binary.Read(reportReader, binary.LittleEndian, r)
	if err != nil {
		panic(err)
	}
}
func (r *Input0x20) Marshal() []byte {
	r.CRC = ReportCRC(PS_FEATURE_CRC32_SEED, r)

	buff := bytes.Buffer{}
	if err := binary.Write(&buff, binary.LittleEndian, r); err != nil {
		panic(err)
	}

	// integrity check
	if buff.Len() != DS_FEATURE_REPORT_FIRMWARE_INFO_SIZE {
		panic(fmt.Sprintf("[%T] len(%d) != %d", r, buff.Len(), DS_FEATURE_REPORT_FIRMWARE_INFO_SIZE))
	}

	return buff.Bytes()
}

// handle0x20 handles HW/FW info (bt only)
func (d *Device) handle0x20(data []byte) {
	if !ReportCRCIsValid(PS_FEATURE_CRC32_SEED, data) {
		fmt.Printf("ERROR handle0x20 CRC check failed [% X]\n", data)
	}

	r := Input0x20{}
	r.Unmarshal(data)

	d.Info.SetHW(r.HardwareVersion, r.FirmwareVersion)

	//fmt.Printf("[%T] %#v\n", r, r)
}
