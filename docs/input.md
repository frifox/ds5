## Buttons
    // ds5.Device.Buttons
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

Buttons have `KeyDown` / `KeyUp` callbacks. Ex:

    x := &dev.Buttons.Cross
	x.OnKeyDown = func() {
		fmt.Println("X pressed")
	}
	x.OnKeyUp = func() {
		fmt.Println("X released")
	}

## Joysticks / Throttles
    // ds5.Device.Axis
    type Axis struct {
        Left  Joystick
        Right Joystick
    
        L2 Throttle
        R2 Throttle
    }

Axis have `OnChange` callbacks. Ex:

    dev.Axis.Left.OnChange = func(x, y float64) {
        fmt.Printf("Left Joystic X:%.3f Y:%.3f\n", x, y)
    }
    dev.Axis.Right.OnChange = func(x, y float64) {
        fmt.Printf("Right Joystic X:%.3f Y:%.3f\n", x, y)
    }

## Touchpad
Track 1 or 2 finger touches across a 1920x1080 touchpad

    // ds5.Device.Touchpad
    type Touchpad struct {
        Touch1 Touch
        Touch2 Touch
    }

Touches have `OnActive` / `OnInactive` callbacks. Ex:

    t1 := &dev.Touchpad.Touch1
	t1.OnActive = func(id uint8, x int, y int) {
		fmt.Printf("Touch1 Active [ID:%d X:%d Y:%d]\n", id, x, y)
	}
	t1.OnInactive = func(id uint8) {
		fmt.Printf("Touch1 Inactive [ID:%d]\n", id)
	}

## Battery

Monitor controller battery status and whether it's currently charging or not.

    // dev.Battery
    type Battery struct {
        Percent  uint8
        Status   string
    }

You can monitor changes (see [src](https://github.com/frifox/ds5/blob/master/handle_0x31.go#L167) for details) via `OnChange` callback, ex:

    dev.Battery.OnChange = func() {
		fmt.Printf("Battery is %s (%d%%)\n", dev.Battery.Status, dev.Battery.Percent)
	}

## Bus

    // dev.Bus
    type Bus struct {
        Type string
    }

Data packets over BT are CRC32 signed and packets over USB are not. Bus value is set to `bt` / `usb` every time `dev.Watch()` is called, where packet crc is checked for the first time.

    dev.Bus.OnChange = func() {
		fmt.Printf("DS5 is now connected via %s\n", dev.Bus.Type)
	}
    

## AliveFor

As reported by the controller. For less verbosity I round this event value to the nearest second.

    // dev.AliveFor
    type AliveFor struct {
        Duration time.Duration
    }

Bind to OnChange event ex:

    dev.AliveFor.OnChange = func(t time.Duration) {
    	fmt.Printf("AliveFor %s\n", t.String())
    }