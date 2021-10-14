
    type Device struct {
        //

        LightBar   LightBar
        PlayerLEDs PlayerLEDs
        Mic        Mic
        Rumble     Rumble

        //
    }

    // this doc will refer to Device as `dev`
    var dev ds5.Device

## RGB LightBar
Full RGB lightbar around the touchpad.

    type LightBar struct {
        Red   uint8
        Green uint8
        Blue  uint8
    }

Pass appropriate property to device.ApplyProps() to apply new state on the controller. Ex:

    // set new state
    dev.LightBar = ds5.LightBar{
        Red: 0,
        Green: 255,
        Bluw: 0,
    ]
    dev.ApplyProps()

    // or modify curent one (will leave untouched values as is)
    dev.LightBar.Green = 255
    dev.ApplyProps()

    // or use a preset color
    dev.LightBar.SetGreen()
    dev.ApplyProps()

## Player LEDs
5 white LEDs below the touchpad

    type PlayerLEDs [5]bool

Control individual LEDs:

    dev.PlayerLEDs[0] = true
    dev.PlayerLEDs[4] = true
    dev.ApplyProps()

Or, light up leds according to a predefined PlayerID/VolumeBar map ([see src](https://github.com/frifox/ds5/blob/master/device_leds.go))

    // mimic PS5 Player identification 
    dev.PlayerLEDs.SetPlayer(2)
    dev.ApplyProps()

    // or light up leds in a row from the left
    dev.PlayerLEDs.SetBar(3)
    dev.ApplyProps()

## Rumble

Control left / right rumble motors. Left motor is slow/deep and right motor is fast/light.

    type Rumble struct {
        Left  uint8
        Right uint8
    }

Rumble left motor at max force for 1 sec, ex:

    dev.Rumble.Left = 255
    dev.ApplyProps()

    time.Sleep(time.Second)

    dev.Rumble.Left = 0
    dev.ApplyProps()

## Mic

Control the Mute button LED and the onboard Mic state individually.

    type Mic struct {
        LED   bool
        Muted bool
    }

Turn on LED but keep Mic on, ex:

    dev.Mic.LED = true
    dev.Mic.Muted = false
    dev.ApplyProps()