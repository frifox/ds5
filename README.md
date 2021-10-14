# DualSense in Go

Heavily based on the official [hid-playstation](https://github.com/torvalds/linux/blob/master/drivers/hid/hid-playstation.c) linux kernel driver.

## Linux / MacOS

This package uses [sstallion/go-hid](https://github.com/sstallion/go-hid) which provides Go bindings for [signal11/hidapi](https://github.com/signal11/hidapi).

[signal11/hidapi](https://github.com/signal11/hidapi) is a multi-platform C library for interfacing with USB/BT HID-Class devices.

PS: Linux: `apt-get install libhidapi-dev libudev-dev`

PS: MacOS: `brew install hidapi`
# Usage
For working examples, check [/example](https://github.com/frifox/ds5/tree/master/example) folder.
    
    package main
    
    import (
        "fmt"
        "github.com/frifox/ds5"
    )
    
    func main() {
        dev := ds5.Device{}
        if err := dev.Find(); err != nil {
            fmt.Printf("Couldn't find DS5: %v\n", err)
            return
        }
    
        dev.Buttons.Cross.OnKeyDown = func() {
            fmt.Printf("X pressed!\n")
        }
    
        fmt.Printf("Watching DS5 for events")
        dev.Watch()
        fmt.Printf("DS5 disappeared\n")
    }

# DS5: ds5.Device{}
This is the main device struct.

Use `Buttons`, `Axis`, `Motion` properties for input events.

Use `Bus`, `Battery`, `AliveFor` for status info.

Use `LightBar`, `PlayerLED`, and `Mute` for controlling the controller hardware.

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
    type Bus string

Data packets over BT are CRC32 signed and packets over USB are not. Bus value is set to `bt` / `usb` every time `dev.Watch()` is called, where packet crc is checked for the first time.

    // TODO add OnChange callback for Bus

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

## RGB LightBar
Full RGB lightbar around the touchpad.

    // ds5.Device.LightBar
    type LightBar struct {
        Red   uint8
        Green uint8
        Blue  uint8
    }

Pass appropriate property to device.ApplyProps() to apply new state on the controller. Ex:

    // dev := ds5.Device{}

    // set new state
    rgb := ds5.LightBar{
        Red: 0,
        Green: 255,
        Bluw: 0,
    ]
    dev.ApplyProps(rgb)

    // or modify curent one (will leave untouched values as is)
    dev.LightBar.Green = 255
    dev.ApplyProps(dev.Lightbar)

## Player LEDs
5 white LEDs below the touchpad

    // ds5.PlayerLEDs
    type PlayerLEDs [5]bool

Control individual LEDs:
    
    dev.PlayerLEDs[0] = true
    dev.PlayerLEDs[4] = true
    dev.ApplyProps(dev.PlayerLEDs)

Or, light up leds according to a predefined PlayerID/VolumeBar map ([see src](https://github.com/frifox/ds5/blob/master/device_leds.go))

    // mimic PS5 Player identification 
    dev.PlayerLEDs.SetPlayer(2)
    dev.ApplyProps(dev.PlayerLEDs)

    // or light up leds in a row from the left
    dev.PlayerLEDs.SetBar(3)
    dev.ApplyProps(dev.PlayerLEDs)

## Rumble

Control left / right rumble motors. Left motor is slow/deep and right motor is fast/light.

    // dev.Rumble
    type Rumble struct {
        Left  uint8
        Right uint8
    }

Rumble left motor for 1 sec, ex:

    dev.Rumble.Left = 255
    dev.ApplyProps(dev.Rumble)

    time.Sleep(time.Second)

    dev.Rumble.Left = 0
    dev.ApplyProps(dev.Rumble)

## Mic

Control the Mute button LED and the onboard Mic state individually.

    // ds5.Mic
    type Mic struct {
        LED   bool
        Muted bool
    }

Turn on LED but keep Mic on, ex:

    dev.Mic.LED = true
    dev.Mic.Muted = false
    dev.ApplyProps(dev.Mic)