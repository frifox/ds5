package ds5

import (
	"fmt"
	"time"
)

type Bus string

func (b *Bus) Set(value string) {
	if value != "bt" && value != "usb" {
		panic("bus.Set() accepts only 'bt' and 'usb'")
	}

	*b = Bus(value)

	fmt.Printf("[BUS] %s\n", *b)
}

type AliveFor struct {
	Duration time.Duration
	OnChange func(time.Duration)
}

func (t *AliveFor) Set(duration time.Duration) {
	if t.Duration == duration {
		return
	}

	t.Duration = duration

	if t.OnChange != nil {
		t.OnChange(t.Duration)
	}

}

type Battery struct {
	Percent  uint8
	Status   string
	OnChange func()
}

func (b *Battery) Set(status string, percent uint8) {
	if b.Status == status && b.Percent == percent {
		return // nothing changed
	}

	b.Status = status
	b.Percent = percent

	// any callbacks?
	if b.OnChange != nil {
		b.OnChange()
	}
}

type Mic struct {
	LED   bool
	Muted bool
}

type Rumble struct {
	Left  uint8 // slow rumble on left
	Right uint8 // fast rumble on right
}

type OutputSequencer uint8

func (s *OutputSequencer) Get() uint8 {
	seq := *s // return current seq

	*s++          // next seq
	*s = *s & 0xf // truncate it to uint4 (0 to 15)

	return uint8(seq)
}
