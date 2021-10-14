package ds5

import (
	"context"
	"fmt"
	"time"
)

type Buttons struct {
	// right buttons
	Square   Button
	Cross    Button
	Circle   Button
	Triangle Button

	// left buttons
	DPadUp    Button
	DPadRight Button
	DPadDown  Button
	DPadLeft  Button

	// center buttons
	Share    Button
	Touchpad Button
	Options  Button
	PS       Button
	Mute     Button

	// back of controller
	L1 Button
	R1 Button

	L2 Button
	R2 Button

	// joysticks
	Left  Button
	Right Button
}

type Button struct {
	Pressed   bool
	OnKeyDown func()
	OnKeyUp   func()

	LongPressTimeout time.Duration
	OnLongPress      func()

	context.Context
	cancel context.CancelFunc
}

func (b *Button) Set(state bool) {
	if b.Pressed == state {
		return // nothing changed
	}

	b.Pressed = state

	// generic callbacks?
	if b.Pressed && b.OnKeyDown != nil {
		go b.OnKeyDown()
	}
	if !b.Pressed && b.OnKeyUp != nil {
		go b.OnKeyUp()
	}

	// support long press callbacks
	if b.OnLongPress != nil {
		if b.LongPressTimeout == 0 {
			b.LongPressTimeout = time.Second
		}

		// KeyDown
		if b.Pressed {
			b.Context, b.cancel = context.WithTimeout(context.Background(), b.LongPressTimeout)

			go func() {
				<-b.Done()

				switch b.Err() {
				case context.DeadlineExceeded:
					go b.OnLongPress()
				case context.Canceled:
					// KeyUp happened before timeout
				default:
					fmt.Printf("LongPress unhandled exit reason: %v\n", b.Err())
				}
			}()
		}

		// KeyUp
		if !b.Pressed {
			b.cancel()
		}
	}
}
