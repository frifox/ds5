package main

import (
	"fmt"
	"github.com/frifox/ds5"
	"math"
	"sync"
	"time"
)

func main() {
	ds5.PrintAllHIDs()

	dev := &ds5.Device{}
	dev.LightBar.Set(255, 255, 255)

	// there are ~128 steps from center to 100% left/right/up/down
	dev.Axis.Left.DeadZone = 10.0 / 128  // ignore -10 to +10 from center
	dev.Axis.Right.DeadZone = 10.0 / 128 // ignore -10 to +10 from center

	// I like when up is +1 and down is -1
	dev.Axis.Left.InvertY = true
	dev.Axis.Right.InvertY = true

	dev.PlayerLEDs[1] = false

	// bind to some events before we start ds5 watcher
	setButtonCallbacks(dev)
	setAxisCallbacks(dev)
	setTouchCallbacks(dev)
	setMiscCallbacks(dev)

	// find ds5, connect, watch. Loop again if disconnected
	for {
		fmt.Printf("Looking for DS5\n")
		for {
			_ = dev.Find()
			if dev.Found() {
				break
			}
			time.Sleep(time.Second)
		}

		fmt.Printf("Starting DS5 watcher\n")
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			dev.Watch()
			wg.Done()
		}()

		// we can set callbacks here too, after watcher has started
		time.Sleep(time.Second)
		setAdvancedCallbacks(dev)

		// wait for Watch goroutine to die
		wg.Wait()
		fmt.Printf("DS5 watcher died. Disconnected?\n")
	}
}

func setMiscCallbacks(dev *ds5.Device) {
	dev.Battery.OnChange = func() {
		fmt.Printf("[Battery] %s (%d%%)\n", dev.Battery.Status, dev.Battery.Percent)
	}

	dev.AliveFor.OnChange = func(t time.Duration) {
		fmt.Printf("[AliveFor] %s\n", t.String())
	}
}

func setButtonCallbacks(dev *ds5.Device) {
	// simple KeyDown / KeyUp callbacks
	setDownUp := func(name string, b *ds5.Button) {
		b.OnKeyDown = func() {
			fmt.Printf("[%s] KeyDown\n", name)
		}
		b.OnKeyUp = func() {
			fmt.Printf("[%s] KeyUp\n", name)
		}
	}

	setDownUp("DPadUp", &dev.Buttons.DPadUp)
	setDownUp("DPadRight", &dev.Buttons.DPadRight)
	setDownUp("DPadDown", &dev.Buttons.DPadDown)
	setDownUp("DPadLeft", &dev.Buttons.DPadLeft)

	setDownUp("Touchpad", &dev.Buttons.Touchpad)

	setDownUp("Left", &dev.Buttons.Left)
	setDownUp("Right", &dev.Buttons.Right)

	// toggle mute button
	dev.Buttons.Mute.OnKeyDown = func() {
		dev.Mic.LED = !dev.Mic.LED  // toggle LED
		dev.Mic.Muted = dev.Mic.LED // and match mic to the LED
		dev.ApplyProps(dev.Mic)
	}

	// switch LightBar color
	dev.Buttons.Square.OnKeyDown = func() {
		dev.ApplyProps(ds5.LightBar{255, 0, 0})
	}
	dev.Buttons.Cross.OnKeyDown = func() {
		dev.ApplyProps(ds5.LightBar{0, 255, 0})
	}
	dev.Buttons.Circle.OnKeyDown = func() {
		dev.ApplyProps(ds5.LightBar{0, 0, 255})
	}
	dev.Buttons.Triangle.OnKeyDown = func() {
		dev.ApplyProps(ds5.LightBar{255, 255, 255})
	}

	// rumble Left while holding
	dev.Buttons.L1.OnKeyDown = func() {
		dev.Rumble.Left = 255
		dev.ApplyProps(dev.Rumble)
	}
	dev.Buttons.L1.OnKeyUp = func() {
		dev.Rumble.Left = 0
		dev.ApplyProps(dev.Rumble)
	}

	// rumble Right while holding
	dev.Buttons.R1.OnKeyDown = func() {
		dev.Rumble.Right = 255
		dev.ApplyProps(dev.Rumble)
	}
	dev.Buttons.R1.OnKeyUp = func() {
		dev.Rumble.Right = 0
		dev.ApplyProps(dev.Rumble)
	}
}

func setAdvancedCallbacks(dev *ds5.Device) {
	// button hold / release
	dev.Buttons.PS.Release = make(chan bool)
	dev.Buttons.PS.OnKeyDown = func() {
		fmt.Printf("[PS] Starting hold\n")
		go func() {
			for {
				tick := time.NewTicker(time.Millisecond * 100)
				select {
				case <-tick.C:
					fmt.Printf("[PS] Holding\n")
				case <-dev.Buttons.PS.Release:
					fmt.Printf("[PS] Done holding\n")
					return
				}
			}
		}()
	}
	dev.Buttons.PS.OnKeyUp = func() {
		dev.Buttons.PS.Release <- true
	}

	// volume bar
	var currentBar uint8
	dev.Buttons.Share.OnKeyDown = func() {
		currentBar--
		if currentBar > 5 {
			currentBar = 5
		}
		dev.PlayerLEDs.SetBar(currentBar)
		dev.ApplyProps(dev.PlayerLEDs)
	}
	dev.Buttons.Options.OnKeyDown = func() {
		currentBar++
		if currentBar > 5 {
			currentBar = 0
		}
		dev.PlayerLEDs.SetBar(currentBar)
		dev.ApplyProps(dev.PlayerLEDs)
	}
}

func setAxisCallbacks(dev *ds5.Device) {
	dev.Axis.Left.OnChange = func(x float64, y float64) {
		fmt.Printf("[Left] X:%.3f Y:%.3f\n", x, y)
	}
	dev.Axis.Right.OnChange = func(x float64, y float64) {
		fmt.Printf("[Right] X:%.3f Y:%.3f\n", x, y)
	}

	dev.Axis.L2.OnChange = func(z float64) {
		fmt.Printf("[L2] Z:%.3f\n", z)
	}
	dev.Axis.R2.OnChange = func(z float64) {
		fmt.Printf("[R2] Z:%.3f\n", z)
	}
}
func setTouchCallbacks(dev *ds5.Device) {
	t1 := &dev.Touchpad.Touch1
	t1.OnActive = func(id uint8, x int, y int) {
		type XY struct {
			X float64
			Y float64
		}
		type Color struct {
			Home XY
			Far  XY
			Max  float64
		}
		max := func(home XY, far XY) float64 {
			return math.Sqrt(math.Pow(home.X-far.X, 2) + math.Pow(home.Y-far.Y, 2))
		}

		var r, g, b Color

		r.Home = XY{0, 0}
		r.Far = XY{1920, 1080}
		r.Max = max(r.Home, r.Far)

		g.Home = XY{960, 1080}
		g.Far = XY{1920, 0}
		g.Max = max(g.Home, g.Far)

		b.Home = XY{1920, 0}
		b.Far = XY{0, 1080}
		b.Max = max(b.Home, b.Far)

		R := ds5.ConvertRange(0, r.Max, 255, 0, t1.DistanceTo(r.Home.X, r.Home.Y))
		G := ds5.ConvertRange(0, g.Max, 255, 0, t1.DistanceTo(g.Home.X, g.Home.Y))
		B := ds5.ConvertRange(0, b.Max, 255, 0, t1.DistanceTo(b.Home.X, b.Home.Y))

		dev.LightBar.Red = uint8(R)
		dev.LightBar.Green = uint8(G)
		dev.LightBar.Blue = uint8(B)

		dev.ApplyProps(dev.LightBar)
		//fmt.Printf("[Touch1] ID:%d X:%d Y:%d\n", id, x, y)
	}

	t2 := &dev.Touchpad.Touch2
	t2.OnActive = func(id uint8, x int, y int) {
		fmt.Printf("[Touch2] ID:%d X:%d Y:%d\n", id, x, y)
	}
	t2.OnInactive = func(id uint8) {
		fmt.Printf("[Touch2] ID:%d Inactive\n", id)
	}
}
