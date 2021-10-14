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

func (d *Device) emit0x2(extra ...interface{}) {
	r := Output0x2{
		ReportID: DS_OUTPUT_REPORT_USB,
	}

	props := []interface{}{
		d.PlayerLEDs,
		d.LightBar,
		d.Rumble,
		d.Mic,
	}
	props = append(props, extra...)

	for _, prop := range props {
		r.ApplyProp(prop)
	}

	d.writer <- &r
}
