package ds5

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sstallion/go-hid"
	"hash/crc32"
	"math"
	"runtime"
	"strconv"
)

func PrintAllHIDs() {
	err := hid.Enumerate(USB_VENDOR_ID_SONY, USB_DEVICE_ID_SONY_PS5_CONTROLLER, EnumPrinter)
	if err != nil {
		fmt.Printf("ERROR hid.Enumerate(): %v\n", err)
	}
}

func EnumPrinter(i *hid.DeviceInfo) error {
	fmt.Printf("Device Information:\n")
	fmt.Printf("\tPath         %s\n", i.Path)
	fmt.Printf("\tVendorID     %#04x\n", i.VendorID)
	fmt.Printf("\tProductID    %#04x\n", i.ProductID)
	fmt.Printf("\tSerialNbr    %s\n", i.SerialNbr)
	fmt.Printf("\tReleaseNbr   %x.%x\n", i.ReleaseNbr>>8, i.ReleaseNbr&0xff)
	fmt.Printf("\tMfrStr       %s\n", i.MfrStr)
	fmt.Printf("\tProductStr   %s\n", i.ProductStr)
	fmt.Printf("\tUsagePage    %#x\n", i.UsagePage)
	fmt.Printf("\tUsage        %#x\n", i.Usage)
	fmt.Printf("\tInterfaceNbr %d\n", i.InterfaceNbr)

	return nil
}

func ConvertRange(inputMin, inputMax, outputMin, outputMax float64, inputValue interface{}, clip ...bool) (outputValue float64) {
	inputValueAsFloat := 0.0
	switch i := inputValue.(type) {
	case uint8:
		inputValueAsFloat = float64(i)
	case int:
		inputValueAsFloat = float64(i)
	case int64:
		inputValueAsFloat = float64(i)
	case float64:
		inputValueAsFloat = i
	default:
		fmt.Printf("ERR unknown type to convert to float64: %T\n", inputValue)
		return 0
	}

	// how far does input range spread?
	inputRange := inputMax - inputMin

	// and how far along that range are we?
	inputPercent := (inputValueAsFloat - inputMin) / inputRange

	// how far does our output range spread?
	outputRange := outputMax - outputMin

	// how far along output range are we?
	outputValue = outputMin + outputRange*inputPercent

	// should we clip out-of-range values?
	if len(clip) > 0 && clip[0] == true {
		if outputValue < outputMin {
			outputValue = outputMin
		}
		if outputValue > outputMax {
			outputValue = outputMax
		}
	}

	return
}
func RemoveDeadZone(deadZone, value float64) float64 {
	deadZone = math.Abs(deadZone)
	if math.Abs(value) <= deadZone {
		return 0
	} else {
		if value > 0 {
			return ConvertRange(deadZone, 1, 0, 1, value)
		} else {
			return ConvertRange(-deadZone, -1, 0, -1, value)
		}
	}
}

// ref: https://github.com/torvalds/linux/blob/master/drivers/hid/hid-playstation.c# ref:L445
func ReportCRCIsValid(seed uint8, report []byte) bool {
	reportData := report[:len(report)-4]
	reportCRC := report[len(report)-4:]

	// different seed bytes are prepended to different reports
	checkData := append([]byte{seed}, reportData...)
	checkCRC := make([]byte, 4)
	binary.LittleEndian.PutUint32(checkCRC, crc32.ChecksumIEEE(checkData))

	if !bytes.Equal(reportCRC, checkCRC) {
		return false
	}

	return true
}
func ReportCRC(seed uint8, report interface{}) uint32 {
	buff := bytes.Buffer{}
	if err := binary.Write(&buff, binary.LittleEndian, report); err != nil {
		panic(err)
	}

	data := buff.Bytes()[:buff.Len()-4]
	seededData := append([]byte{seed}, data...)

	return crc32.ChecksumIEEE(seededData)
}

func goID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
