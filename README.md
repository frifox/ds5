# DualSense in Go

Heavily based on the official [hid-playstation](https://github.com/torvalds/linux/blob/master/drivers/hid/hid-playstation.c) linux kernel driver.

## Linux / MacOS

* This package uses [sstallion/go-hid](https://github.com/sstallion/go-hid) which provides Go bindings for [signal11/hidapi](https://github.com/signal11/hidapi).
* [signal11/hidapi](https://github.com/signal11/hidapi) is a multi-platform C library for interfacing with USB/BT HID-Class devices.
* Linux: `apt-get install libhidapi-dev libudev-dev`
* MacOS: `brew install hidapi`
  * on arm64 you may need to manually add path, found via `brew ls hidapi`,
  * ie: `CGO_CFLAGS="-I /opt/homebrew/Cellar/hidapi/0.11.0/include" CGO_LDFLAGS="-L /opt/homebrew/Cellar/hidapi/0.11.0/lib/" go build .`)
* Windows: I don't have one to test...

# Usage
For working examples, check [examples](https://github.com/frifox/ds5/tree/master/examples) folder.
    
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
    
        dev.Buttons.Square.OnKeyDown = func() {
            fmt.Printf("Square pressed\n")
    
            dev.LightBar.SetRed()
            dev.ApplyProps()
        }
        dev.Buttons.Square.OnKeyUp = func() {
            fmt.Printf("Square released\n")
    
            dev.LightBar.SetGreen()
            dev.ApplyProps()
        }
    
        fmt.Printf("Watching DS5 for events\n")
        go dev.Run()
        
        <- dev.Done()
        fmt.Printf("DS5 disappeared\n")
    }


# More info?
See [docs/input.md](https://github.com/frifox/ds5/tree/master/docs/input.md) & [docs/output.md](https://github.com/frifox/ds5/tree/master/docs/output.md)