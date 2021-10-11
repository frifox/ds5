package main

import (
	"fmt"
	"github.com/frifox/ps5"
)

func main() {
	dev := ds5.Device{}
	if err := dev.Find(); err != nil {
		fmt.Printf("Couldn't find DS5: %v\n", err)
		return
	}

	setCallbacks(&dev)

	fmt.Printf("Watching DS5 for events")
	dev.Watch()
	fmt.Printf("DS5 disappeared\n")
}

func setCallbacks(dev *ds5.Device) {
	dev.Buttons.Cross.OnKeyDown = func() {
		fmt.Printf("X pressed!\n")
	}

}
