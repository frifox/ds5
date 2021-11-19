package ds5

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

const DS_FEATURE_REPORT_PAIRING_INFO = 0x09
const DS_FEATURE_REPORT_PAIRING_INFO_SIZE = 20

type Input0x9 struct {
	ReportID byte

	MAC [6]byte

	Unknown [13]byte
}

func (r *Input0x9) Unmarshal(data []byte) {
	reportReader := bytes.NewReader(data)
	err := binary.Read(reportReader, binary.LittleEndian, r)

	// MAC comes in reversed. Fix that
	for i := 0; i < len(r.MAC)/2; i++ {
		j := len(r.MAC) - i - 1
		r.MAC[i], r.MAC[j] = r.MAC[j], r.MAC[i]
	}

	if err != nil {
		panic(err)
	}
}

// handle0x9 handles MAC info
func (d *Device) handle0x9(data []byte) {
	if ReportCRCIsValid(PS_FEATURE_CRC32_SEED, data) {
		d.Bus.Set("bt")
	} else {
		d.Bus.Set("usb")
	}

	r := Input0x9{}
	r.Unmarshal(data)

	mac := fmt.Sprintf("% X", r.MAC)
	mac = strings.ReplaceAll(mac, " ", ":")
	d.MAC.Set(mac)

	//fmt.Printf("[%T] %#v\n", r, r)
}
