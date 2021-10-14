package ds5

// Axis values -1 to 1
type Axis struct {
	Left  Joystick
	Right Joystick

	L2 Throttle
	R2 Throttle
}

type Joystick struct {
	DeadZone float64
	InvertX  bool
	InvertY  bool

	X        float64
	Y        float64
	OnChange func(float64, float64)
}

func (j *Joystick) Set(x uint8, y uint8) {
	X := ConvertRange(0, 255, -1, 1, x)
	Y := ConvertRange(0, 255, -1, 1, y)

	X = RemoveDeadZone(j.DeadZone, X)
	Y = RemoveDeadZone(j.DeadZone, Y)

	if j.X == X && j.Y == Y {
		return // nothing changed
	}

	// flip axis?
	if j.InvertX {
		X = -X
	}
	if j.InvertY {
		Y = -Y
	}

	j.X = X
	j.Y = Y

	// any callbacks?
	if j.OnChange != nil {
		go j.OnChange(j.X, j.Y)
	}
}

type Throttle struct {
	Z        float64
	OnChange func(float64)
}

func (t *Throttle) Set(z uint8) {
	Z := 0.0

	Z = ConvertRange(0, 255, 0, 1, z)

	if t.Z == Z {
		return // nothing changed
	}

	t.Z = Z

	// any callbacks?
	if t.OnChange != nil {
		go t.OnChange(t.Z)
	}
}
