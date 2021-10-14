package ds5

type Bus struct {
	Type     string
	OnChange func(string)
}

func (b *Bus) Set(value string) {
	if value != "bt" && value != "usb" {
		panic("bus.Set() accepts only 'bt' and 'usb'")
	}

	if b.Type == value {
		return
	}

	b.Type = value

	// callbacks
	if b.OnChange != nil {
		go b.OnChange(b.Type)
	}
}
