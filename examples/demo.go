package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/frifox/ds5"
	"math"
	"time"
)

var dev *ds5.Device

func main() {
	ds5.PrintAllHIDs()

	dev = &ds5.Device{}
	dev.LightBar.Set(0, 255, 0)

	// DS5 axis resolution is uint8: 0 to 255 (~128 is center)
	dev.Axis.Left.DeadZone = 10.0 / 128  // ignore -10 to +10 from center
	dev.Axis.Right.DeadZone = 10.0 / 128 // ignore -10 to +10 from center

	// I like when up is +1 and down is -1
	dev.Axis.Left.InvertY = true
	dev.Axis.Right.InvertY = true

	// bind to events before we start ds5 watcher
	setLeftButtonCallbacks()
	setRightButtonCallbacks()
	setCenterButtonCallbacks()
	setBackButtonsCallbacks()
	setAxisCallbacks()
	setTouchCallbacks()
	setMiscCallbacks()
	// NOTE: you CAN set bind after watcher started too...

	// find ds5, connect, watch... Loop again if disconnected
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
		go dev.Run()

		// we can set callbacks here too, after watcher has started
		time.Sleep(time.Second)

		// wait for Watch goroutine to die
		<-dev.Done()
		fmt.Printf("DS5 watcher died. Disconnected?\n")
	}
}

func setLeftButtonCallbacks() {
	setDownUp := func(name string, b *ds5.Button) {
		b.OnKeyDown = func() {
			fmt.Printf("[%s] KeyDown\n", name)
		}
		b.OnKeyUp = func() {
			fmt.Printf("[%s] KeyUp\n", name)
		}

		// called 1s after KeyDown if held
		b.OnLongPress = func() {
			fmt.Printf("[%s] LongPress\n", name)
		}
	}

	setDownUp("DPadUp", &dev.Buttons.DPadUp)
	setDownUp("DPadRight", &dev.Buttons.DPadRight)
	setDownUp("DPadDown", &dev.Buttons.DPadDown)
	setDownUp("DPadLeft", &dev.Buttons.DPadLeft)

	setDownUp("Touchpad", &dev.Buttons.Touchpad)

	setDownUp("Left", &dev.Buttons.Left)
	setDownUp("Right", &dev.Buttons.Right)
}
func setRightButtonCallbacks() {
	// switch LightBar color
	dev.Buttons.Square.OnKeyDown = func() {
		dev.LightBar.SetRed()
		dev.ApplyProps()
		fmt.Printf("[Square] LightBar.SetRed()\n")
	}
	dev.Buttons.Cross.OnKeyDown = func() {
		dev.LightBar.SetGreen()
		dev.ApplyProps()
		fmt.Printf("[Cross] LightBar.SetGreen()\n")
	}
	dev.Buttons.Circle.OnKeyDown = func() {
		dev.LightBar.SetBlue()
		dev.ApplyProps()
		fmt.Printf("[Circle] LightBar.SetBlue()\n")
	}
	dev.Buttons.Triangle.OnKeyDown = func() {
		dev.LightBar = ds5.LightBar{255, 255, 255}
		dev.ApplyProps()
		fmt.Printf("[Triangle] LightBar(255,255,255)\n")
	}
}

func setCenterButtonCallbacks() {
	// button hold
	var press context.Context
	var cancel context.CancelFunc
	dev.Buttons.PS.OnKeyDown = func() {
		press, cancel = context.WithCancel(context.Background())
		fmt.Printf("[PS] Starting hold loop\n")
		go func() {
			// do something every 100ms until done holding
			for {
				tick := time.NewTicker(time.Millisecond * 100)
				select {
				case <-tick.C:
					fmt.Printf("[PS] Holding\n")
				case <-press.Done():
					fmt.Printf("[PS] Done holding\n")
					return
				}
			}
		}()
	}
	dev.Buttons.PS.OnKeyUp = func() {
		cancel()
	}

	// volume bar +/- via Share/Options keys
	var currentBar uint8
	dev.Buttons.Share.OnKeyDown = func() {
		currentBar--
		if currentBar > 5 {
			currentBar = 5
		}
		dev.PlayerLEDs.SetDot(currentBar)
		dev.ApplyProps()
		fmt.Printf("[Share] SetDot(%d)\n", currentBar)
	}
	dev.Buttons.Options.OnKeyDown = func() {
		currentBar++
		if currentBar > 5 {
			currentBar = 0
		}
		dev.PlayerLEDs.SetDot(currentBar)
		dev.ApplyProps()
		fmt.Printf("[Options] SetDot(%d)\n", currentBar)
	}

	// toggle mute button
	dev.Buttons.Mute.OnKeyDown = func() {
		dev.Mic.LED = !dev.Mic.LED  // toggle LED
		dev.Mic.Muted = dev.Mic.LED // and match mic to the LED
		dev.ApplyProps()
		fmt.Printf("[Mute] Mic.LED(%t)\n", dev.Mic.LED)
	}
}

func setBackButtonsCallbacks() {
	// rumble Left while holding
	dev.Buttons.L1.OnKeyDown = func() {
		dev.Rumble.Left = 255
		dev.ApplyProps()
		fmt.Printf("[L1] Rumble(255)\n")
	}
	dev.Buttons.L1.OnKeyUp = func() {
		dev.Rumble.Left = 0
		dev.ApplyProps()
		fmt.Printf("[L1] Rumble(0)\n")
	}

	// rumble Right while holding
	dev.Buttons.R1.OnKeyDown = func() {
		dev.Rumble.Right = 255
		dev.ApplyProps()
		fmt.Printf("[R1] Rumble(255)\n")
	}
	dev.Buttons.R1.OnKeyUp = func() {
		dev.Rumble.Right = 0
		dev.ApplyProps()
		fmt.Printf("[R1] Rumble(0)\n")
	}
}

func setAxisCallbacks() {
	dev.Axis.Left.OnChange = func(a ds5.Joystick) {
		fmt.Printf("[Left] X:%.3f Y:%.3f\n", a.X, a.Y)
	}
	dev.Axis.Right.OnChange = func(a ds5.Joystick) {
		fmt.Printf("[Right] X:%.3f Y:%.3f\n", a.X, a.Y)
	}

	dev.Axis.L2.OnChange = func(a ds5.Throttle) {
		fmt.Printf("[L2] Z:%.3f\n", a.Z)
	}
	dev.Axis.R2.OnChange = func(a ds5.Throttle) {
		fmt.Printf("[R2] Z:%.3f\n", a.Z)
	}
}
func setTouchCallbacks() {
	// change Lightbar color depending on where you touch
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

	t1 := &dev.Touchpad.Touch1
	t1.OnActive = func(t ds5.Touch) {
		R := ds5.ConvertRange(0, r.Max, 255, 0, t.DistanceTo(r.Home.X, r.Home.Y))
		G := ds5.ConvertRange(0, g.Max, 255, 0, t.DistanceTo(g.Home.X, g.Home.Y))
		B := ds5.ConvertRange(0, b.Max, 255, 0, t.DistanceTo(b.Home.X, b.Home.Y))

		dev.LightBar.Red = uint8(R)
		dev.LightBar.Green = uint8(G)
		dev.LightBar.Blue = uint8(B)

		dev.ApplyProps()
		//fmt.Printf("[Touch1] ID:%d X:%d Y:%d\n", id, x, y)
	}

	t2 := &dev.Touchpad.Touch2
	t2.OnActive = func(t ds5.Touch) {
		fmt.Printf("[Touch2] ID:%d X:%d Y:%d\n", t.ID, t.X, t.Y)
	}
	t2.OnInactive = func(t ds5.Touch) {
		fmt.Printf("[Touch2] ID:%d Inactive\n", t.ID)
	}
}

func setMiscCallbacks() {
	dev.Bus.OnChange = func(b ds5.Bus) {
		fmt.Printf("[Bus] %s\n", b.Type)
	}
	dev.Battery.OnChange = func(b ds5.Battery) {
		fmt.Printf("[Battery] %s (%d%%)\n", b.Status, b.Percent)
	}
	dev.MAC.OnChange = func(i ds5.MAC) {
		fmt.Printf("[MAC] %s\n", i.Address)
	}
	dev.Version.OnChange = func(i ds5.Version) {
		dots := make([]byte, 4)
		binary.BigEndian.PutUint32(dots, i.FirmwareVersion)
		dotted := fmt.Sprintf("%d.%d.%d.%d", dots[0], dots[1], dots[2], dots[3])

		fmt.Printf("[Version] Hardware:0x%x, Firmware:0x%x / %s\n", i.HardwareVersion, i.FirmwareVersion, dotted)
	}

	//dev.AliveFor.OnChange = func(t time.Duration) {
	//	fmt.Printf("[AliveFor] %s\n", t.String())
	//}
}
