package ds5

import "math"

// Touchpad points are across 1920 x 1080 plane
type Touchpad struct {
	Touch1 Touch
	Touch2 Touch
}
type Touch struct {
	ID         uint8
	Active     bool
	X          int
	Y          int
	OnActive   func(Touch)
	OnInactive func(Touch)
}

func (t *Touch) Set(id uint8, active bool, x int, y int) {
	if t.ID == id && t.Active == active && t.X == x && t.Y == y {
		return // nothing changed
	}

	t.Active = active
	t.ID = id
	t.X = x
	t.Y = y

	// any callbacks
	if t.Active && t.OnActive != nil {
		go t.OnActive(*t)
	}
	if !t.Active && t.OnInactive != nil {
		go t.OnInactive(*t)
	}
}
func (t *Touch) DistanceTo(x float64, y float64) float64 {
	first := math.Pow(x-float64(t.X), 2)
	second := math.Pow(y-float64(t.Y), 2)
	return math.Sqrt(first + second)
}
