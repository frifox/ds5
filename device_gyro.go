package ds5

const DS_GYRO_RES_PER_DEG_S = 1024
const DS_GYRO_RANGE = 2048 * DS_GYRO_RES_PER_DEG_S

// TODO normalize Gyro values to something meaningful

type Gyro struct {
	Pitch    float64
	Roll     float64
	Yaw      float64
	OnChange func(Gyro)

	pitchCal GyroCalibration
	rollCal  GyroCalibration
	yawCal   GyroCalibration
}

type GyroCalibration struct {
	Bias int16

	Max int16
	Min int16

	speedPlus  int16
	speedMinus int16

	plus        int16
	minus       int16
	numerator   int16
	denominator int16
}

func (g *Gyro) Set(pitch, yaw, roll float64) {
	if g.Pitch == pitch && g.Yaw == yaw && g.Roll == roll {
		return
	}

	// TODO normalize number to something meaningful
	g.Pitch = pitch // -down, +up
	g.Yaw = -yaw    // -left, +right
	g.Roll = -roll  // -left, +right

	if g.OnChange != nil {
		go g.OnChange(*g)
	}
}
