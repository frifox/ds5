package ds5

import "time"

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
