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

	dev.Buttons.Square.OnKeyDown = func() {
		fmt.Printf("Square pressed\n")

		dev.LightBar.SetRed()
		dev.ApplyProps()
	}
	dev.Buttons.Square.OnKeyUp = func() {
		fmt.Printf("Square released\n")

		dev.LightBar.SetGreen()
		dev.ApplyProps()
	}

	fmt.Printf("Watching DS5 for events\n")
	go dev.Run()

	<-dev.Done()
	fmt.Printf("DS5 disappeared\n")
}
