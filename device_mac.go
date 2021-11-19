package ds5

type MAC struct {
	Address string

	OnChange func(MAC)
}

func (i *MAC) Set(mac string) {
	if i.Address == mac {
		return // nothing changed
	}

	i.Address = mac

	// any callbacks?
	if i.OnChange != nil {
		go i.OnChange(*i)
	}
}
