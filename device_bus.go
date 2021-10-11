package ds5

import "fmt"

type Bus string

func (b *Bus) Set(value string) {
	if value != "bt" && value != "usb" {
		panic("bus.Set() accepts only 'bt' and 'usb'")
	}

	*b = Bus(value)

	fmt.Printf("[BUS] %s\n", *b)
}
