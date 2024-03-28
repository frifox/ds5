```go
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
```
## RGB LightBar
Full RGB lightbar around the touchpad.
```go
type LightBar struct {
    Red   uint8
    Green uint8
    Blue  uint8
}
```
Note: run ApplyProps() to apply values on the controller. If changing multiple props at once, call ApplyProps() after all updates are done, to avoid unnecessarily sending multiple packets to DS5.
```go
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
```
## Player LEDs
5 white LEDs below the touchpad
```go
type PlayerLEDs struct {
    LED [5]bool
    DisableChangeAnimation bool
}
```
Control individual LEDs:
```go
dev.PlayerLEDs[0] = true
dev.PlayerLEDs[4] = true
dev.ApplyProps()
```
Or, light up leds according to a predefined PlayerID/VolumeBar map ([see src](https://github.com/frifox/ds5/blob/master/device_leds.go))
```go
// mimic PS5 Player identification 
dev.PlayerLEDs.SetPlayer(2)
dev.ApplyProps()

// or light up leds in a row from the left
dev.PlayerLEDs.SetBar(3)
dev.ApplyProps()
```
## Rumble

Control left / right rumble motors. Left motor is slow/deep and right motor is fast/light.
```go
type Rumble struct {
    Left  uint8
    Right uint8
}
```
Rumble left motor at max force for 1 sec, ex:
```go
dev.Rumble.Left = 255
dev.ApplyProps()

time.Sleep(time.Second)

dev.Rumble.Left = 0
dev.ApplyProps()
```
## Mic

Control the Mute button LED and the onboard Mic state individually.
```go
type Mic struct {
    LED   bool
    Muted bool
}
```
Turn on LED but keep Mic on, ex:
```go
dev.Mic.LED = true
dev.Mic.Muted = false
dev.ApplyProps()
```