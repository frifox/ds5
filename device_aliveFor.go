package ds5

import "time"

type AliveFor struct {
	Duration time.Duration
	OnChange func(AliveFor)
}

func (t *AliveFor) Set(duration time.Duration) {
	if t.Duration == duration {
		return
	}

	t.Duration = duration

	if t.OnChange != nil {
		go t.OnChange(*t)
	}

}
