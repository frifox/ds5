package ds5

const DS_ACC_RES_PER_G = 8192
const DS_ACC_RANGE = 4 * DS_ACC_RES_PER_G

// TODO use degrees/radians instead of -1/+1?

// Accel is how much gravity is pulling on an axis,
// axis in-line with gravity = 1 or -1
// axis perpendicular with gravity = 0
type Accel struct {
	X        float64 // left to right (ie: Roll)
	Y        float64 // bottom to top (ie: Orientation)
	Z        float64 // front to back (ie: Pitch)
	OnChange func(Accel)

	cal InputAccel
}

func (a *Accel) Set(x, y, z float64) {
	if a.X == x && a.Y == y && a.Z == z {
		return
	}

	a.X = -x // tilt: -left, +right
	a.Y = y  // orientation: -belly-up, +upright
	a.Z = -z // nose: -down, +up

	if a.OnChange != nil {
		go a.OnChange(*a)
	}
}
