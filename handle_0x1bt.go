package ds5

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const Input0x1ReportSize = 10

type Input0x1bt struct {
	ReportID byte

	LeftX, LeftY   byte
	RightX, RightY byte

	Buttons [3]byte

	L2, R2 byte
}

func (r *Input0x1bt) Marshal() []byte {
	buff := bytes.Buffer{}
	if err := binary.Write(&buff, binary.LittleEndian, r); err != nil {
		panic(err)
	}

	// integrity check
	if buff.Len() != Input0x1ReportSize {
		panic(fmt.Sprintf("[%T] len(%d) != %d", r, buff.Len(), Input0x1ReportSize))
	}

	//fmt.Printf("[%T] Marshal len(%d) % X\n", i, buff.Len(), buff.Bytes())

	return buff.Bytes()
}
func (r *Input0x1bt) Unmarshal(data []byte) {
	reportReader := bytes.NewReader(data)
	err := binary.Read(reportReader, binary.LittleEndian, r)
	if err != nil {
		panic(err)
	}
}

func (d *Device) handle0x1bt(report []byte) {
	r := Input0x1bt{}
	r.Unmarshal(report)

	d.Axis.Left.Set(r.LeftX, r.LeftY)
	d.Axis.Right.Set(r.RightX, r.RightY)

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

	// sequence number, 1 to 15?
	if r.Buttons[2]>>2 != 0x0 {
		//b := r.Buttons[2]>>2
		//fmt.Printf("[Input0x1] Buttons[2]>>2 [%08b] [% X] [%d] \n", b, b, b)
	}

	d.Axis.L2.Set(r.L2)
	d.Axis.R2.Set(r.R2)
}
