package ds5

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
		go b.OnChange()
	}
}
