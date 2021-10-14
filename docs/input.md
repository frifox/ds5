##
    type Device struct {
        //

        Bus      Bus
        AliveFor AliveFor
        Battery  Battery
    
        Buttons  Buttons
        Axis     Axis
        Touchpad Touchpad
        Gyro     Gyro
        Accel    Accel

        //
    }

    // this doc will refer to Device as `dev`
    var dev ds5.Device

## Buttons
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
    
Buttons have `KeyDown` / `KeyUp` / `OnLongPress` callbacks. `OnLongPress` will trigger as soon as `LongPressTimeout` is reached after `KeyDown`. If all 3 callbacks are assigned, all 3 can and will trigger at their respective events

    x := &dev.Buttons.Cross

	x.OnKeyDown = func() {
		fmt.Println("X pressed")
	}
	x.OnKeyUp = func() {
		fmt.Println("X released")
	}
    x.OnLongPress = func() {
        fmt.Println("X was held for 1s")
    }

PS: `LongPressTimeout` defaults to 1 sec. You can change it, ex:

    x.LongPressTimeout = time.Duration(time.Millisecond * 500)
    x.OnLongPress = func() {
        fmt.Println("X was held for 500ms")
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

    dev.Axis.L1.OnChange = func(z float64) {
        fmt.Printf("L1 throttle: %.3f\n", z)
    }

## Gyro

How fast you change Pitch, Roll, and Yaw.

    type Gyro struct {
        Pitch    float64
        Roll     float64
        Yaw      float64
    }

Values are -1 to +1.

Pitch: -1 down, +1 up

Yaw: -1 left, +1 right

Roll: -1 left, +1 right

    dev.Gyro.OnChange = func(g *ds5.Gyro) {
        fmt.Printf("Gyroscope: Pitch: %.3f | Roll: %.3f | Yaw: %.3f\n", g.Pitch, g.Roll, g.Yaw)     
    }

## Accel

How much gravity is pulling on an axis. Axis in-line with gravity = -1 / +1, axis perpendicular with gravity = 0;

X: left to right (ie: Roll)

Y: bottom to top (ie: Orientation)

Z: front to back (ie: Pitch)

    
    dev.Accel.OnChange = func(a *ds5.Accel) {
        fmt.Printf("Accelerometer: X: %.3f | Y: %.3f | Z: %.3f\n", a.X, a.Y, a.Z)     
    }

## Touchpad
Track 1 or 2 finger touches across a 1920x1080 touchpad

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

    type Battery struct {
        Percent  uint8
        Status   string
    }

You can monitor changes (see [src](https://github.com/frifox/ds5/blob/master/handle_0x31.go#L167) for details) via `OnChange` callback, ex:

    dev.Battery.OnChange = func() {
		fmt.Printf("Battery is %s (%d%%)\n", dev.Battery.Status, dev.Battery.Percent)
	}

## Bus

    type Bus struct {
        Type string
    }

Data packets over BT are CRC32 signed and packets over USB are not. Bus value is set to `bt` / `usb` every time `dev.Watch()` is called, where packet crc is checked for the first time.

    dev.Bus.OnChange = func() {
		fmt.Printf("DS5 is now connected via %s\n", dev.Bus.Type)
	}
    

## AliveFor

As reported by the controller. For less verbosity I round this event value to the nearest second.

    type AliveFor struct {
        Duration time.Duration
    }

Bind to OnChange event ex:

    dev.AliveFor.OnChange = func(t time.Duration) {
    	fmt.Printf("AliveFor %s\n", t.String())
    }