package ds5

// LightBar values 0 to 255
type LightBar struct {
	Red   uint8
	Green uint8
	Blue  uint8
}
type LightBarInit struct {
	// empty
}

func (l *LightBar) Set(r uint8, g uint8, b uint8) {
	l.Red = r
	l.Green = g
	l.Blue = b
}
func (l *LightBar) SetBlack() {
	l.Set(0, 0, 0)
}

func (l *LightBar) SetWhite() {
	l.Set(255, 255, 255)
}
func (l *LightBar) SetRed() {
	l.Set(255, 0, 0)
}
func (l *LightBar) SetGreen() {
	l.Set(0, 255, 0)
}
func (l *LightBar) SetBlue() {
	l.Set(0, 0, 255)
}

func (l *LightBar) SetYellow() {
	l.Set(255, 255, 0)
}
func (l *LightBar) SetCyan() {
	l.Set(0, 255, 255)
}
func (l *LightBar) SetMagenta() {
	l.Set(255, 0, 255)
}

func (l *LightBar) SetOrange() {
	l.Set(255, 128, 0)
}
func (l *LightBar) SetPurple() {
	l.Set(255, 0, 128)
}
func (l *LightBar) SetAqua() {
	l.Set(0, 255, 128)
}
