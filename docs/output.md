

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