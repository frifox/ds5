package ds5

type Buttons struct {
	Square   Button
	Cross    Button
	Circle   Button
	Triangle Button

	DPadUp    Button
	DPadRight Button
	DPadDown  Button
	DPadLeft  Button

	Share    Button
	Touchpad Button
	Options  Button
	PS       Button
	Mute     Button

	L1 Button
	R1 Button

	L2 Button
	R2 Button

	Left  Button
	Right Button
}

type Button struct {
	Pressed   bool
	OnKeyDown func()
	OnKeyUp   func()

	Release chan bool ``
}

func (b *Button) Set(state bool) {
	if b.Pressed == state {
		return // nothing changed
	}

	b.Pressed = state

	// any callbacks?
	if b.Pressed && b.OnKeyDown != nil {
		b.OnKeyDown()
	}
	if !b.Pressed && b.OnKeyUp != nil {
		b.OnKeyUp()
	}
}
