package ds5

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const DS_OUTPUT_TAG = 0x10

const DS_OUTPUT_REPORT_BT = 0x31
const DS_OUTPUT_REPORT_BT_SIZE = 78
const PS_OUTPUT_CRC32_SEED = 0xA2

type Output0x31 struct {
	ReportID byte //DS_OUTPUT_REPORT_BT 0x31
	SeqTag   byte // (0 to 15)<<4
	Tag      byte // DS_OUTPUT_TAG 0x10

	OutputCommon

	Reserved [24]byte
	CRC      uint32
}

func (r *Output0x31) Marshal() []byte {
	r.CRC = ReportCRC(PS_OUTPUT_CRC32_SEED, r)

	buff := bytes.Buffer{}
	if err := binary.Write(&buff, binary.LittleEndian, r); err != nil {
		panic(err)
	}

	// integrity check
	if buff.Len() != DS_OUTPUT_REPORT_BT_SIZE {
		panic(fmt.Sprintf("[%T] len(%d) != %d", r, buff.Len(), DS_OUTPUT_REPORT_BT_SIZE))
	}
	if !ReportCRCIsValid(PS_OUTPUT_CRC32_SEED, buff.Bytes()) {
		panic("Output0x31 computed CRC integrity failed")
	}

	fmt.Printf("[%T] % X\n", r, buff.Bytes())

	return buff.Bytes()
}

func (d *Device) emit0x31(props ...interface{}) {
	r := Output0x31{
		ReportID: DS_OUTPUT_REPORT_BT,
		SeqTag:   d.OutputSequencer.Get() << 4, // shift seq to 0xf0
		Tag:      DS_OUTPUT_TAG,
	}

	for _, prop := range props {
		r.ApplyProp(prop)
	}

	_, err := d.hid.Write(r.Marshal())
	if err != nil {
		fmt.Printf("[%T] ERR SendFeatureReport | %v |len(%d) [%X]\n", r, err, len(r.Marshal()), r.Marshal())
	} else {
		//fmt.Printf("[Emit0x31] Sent %d bytes\n", n)
	}
}
