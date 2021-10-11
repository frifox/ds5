package main

import (
	"fmt"
	"github.com/frifox/ds5"
)

func main() {
	dev := ds5.Device{}
	if err := dev.Find(); err != nil {
		fmt.Printf("Couldn't find DS5: %v\n", err)
		return
	}

	dev.Buttons.Cross.OnKeyDown = func() {
		fmt.Printf("X pressed\n")
	}
	dev.Buttons.Cross.OnKeyUp = func() {
		fmt.Printf("X released\n")
	}
	dev.Axis.Left.OnChange = func(x, y float64) {
		fmt.Printf("Left Joystic X:%.3f Y:%.3f\n", x, y)
	}

	fmt.Printf("Watching DS5 for events")
	dev.Watch()
	fmt.Printf("DS5 disappeared\n")
}
