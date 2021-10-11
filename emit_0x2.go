package ds5

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const DS_OUTPUT_REPORT_USB = 0x2
const DS_OUTPUT_REPORT_USB_SIZE = 63

type Output0x2 struct {
	ReportID byte //DS_OUTPUT_REPORT_USB 0x2

	OutputCommon

	Reserved [15]byte
}

func (r *Output0x2) Marshal() []byte {
	buff := bytes.Buffer{}
	if err := binary.Write(&buff, binary.LittleEndian, r); err != nil {
		panic(err)
	}

	// integrity check
	if buff.Len() != DS_OUTPUT_REPORT_USB_SIZE {
		panic(fmt.Sprintf("[%T] len(%d) != %d", r, buff.Len(), DS_OUTPUT_REPORT_USB_SIZE))
	}

	//fmt.Printf("[Output0x2] len(%d) [% X]\n", buff.Len(), buff.Bytes())

	return buff.Bytes()
}

func (d *Device) emit0x2(props ...interface{}) {
	r := Output0x2{
		ReportID: DS_OUTPUT_REPORT_USB,
	}

	for _, prop := range props {
		r.ApplyProp(prop)
	}

	_, err := d.hid.Write(r.Marshal())
	if err != nil {
		fmt.Printf("[%T] ERR SendFeatureReport | %v |len(%d) [%X]\n", r, err, len(r.Marshal()), r.Marshal())
	} else {
		//fmt.Printf("[Emit0x2] Sent %d bytes\n", n)
	}
}
