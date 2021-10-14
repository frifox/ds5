package ds5

const DS_OUTPUT_VALID_FLAG0_COMPATIBLE_VIBRATION = 1 << 0
const DS_OUTPUT_VALID_FLAG0_HAPTICS_SELECT = 1 << 1

const DS_OUTPUT_VALID_FLAG1_MIC_MUTE_LED_CONTROL_ENABLE = 1 << 0
const DS_OUTPUT_VALID_FLAG1_POWER_SAVE_CONTROL_ENABLE = 1 << 1
const DS_OUTPUT_POWER_SAVE_CONTROL_MIC_MUTE uint8 = 1 << 4

const DS_OUTPUT_VALID_FLAG1_LIGHTBAR_CONTROL_ENABLE = 1 << 2
const DS_OUTPUT_VALID_FLAG1_RELEASE_LEDS = 1 << 3 // hid-playstation.c doesn't use this constant?
const DS_OUTPUT_VALID_FLAG1_PLAYER_INDICATOR_CONTROL_ENABLE = 1 << 4

const DS_OUTPUT_VALID_FLAG2_LIGHTBAR_SETUP_CONTROL_ENABLE = 1 << 1
const DS_OUTPUT_LIGHTBAR_SETUP_LIGHT_OUT = 1 << 1

type outputSequencer uint8

func (s *outputSequencer) Get() uint8 {
	seq := *s     // return current seq
	*s++          // next seq
	*s = *s & 0xf // truncate it to uint4 (0 to 15)
	return uint8(seq)
}

type OutputCommon struct {
	ValidFlag0 byte
	ValidFlag1 byte

	MotorRight byte
	MotorLeft  byte

	Reserved         [4]byte
	MuteButtonLED    byte
	PowerSaveControl byte
	Reserved2        [28]byte

	ValidFlag2    byte
	Reserved3     [2]byte
	LightBarSetup byte
	LedBrightness byte
	PlayerLEDs    byte

	LightBarRed   byte
	LightBarGreen byte
	LightBarBlue  byte
}

func (r *OutputCommon) ApplyProp(prop interface{}) {
	switch p := prop.(type) {
	// flag0 bit 0, 1
	case Rumble:
		//fmt.Printf("[Emit0x31] %#v\n", p)

		r.ValidFlag0 |= DS_OUTPUT_VALID_FLAG0_HAPTICS_SELECT
		r.ValidFlag0 |= DS_OUTPUT_VALID_FLAG0_COMPATIBLE_VIBRATION

		r.MotorLeft = p.Left
		r.MotorRight = p.Right

	// flag0 bit 2 - 7 ?

	// flag1 bit 0, 1
	case Mic:
		//fmt.Printf("[Emit0x31] %#v\n", p)

		r.ValidFlag1 |= DS_OUTPUT_VALID_FLAG1_MIC_MUTE_LED_CONTROL_ENABLE

		// enable/disable mute LED
		if p.LED {
			r.MuteButtonLED = 0x1
		} else {
			r.MuteButtonLED = 0x0
		}

		r.ValidFlag1 |= DS_OUTPUT_VALID_FLAG1_POWER_SAVE_CONTROL_ENABLE
		if p.Muted {
			r.PowerSaveControl |= DS_OUTPUT_POWER_SAVE_CONTROL_MIC_MUTE
		} else {
			r.PowerSaveControl &= ^DS_OUTPUT_POWER_SAVE_CONTROL_MIC_MUTE
		}

	// flag1 bit 2
	case LightBar:
		//fmt.Printf("[Emit0x31] %#v\n", p)

		r.ValidFlag1 |= DS_OUTPUT_VALID_FLAG1_LIGHTBAR_CONTROL_ENABLE

		// convert 0-255 to 0-128
		r.LightBarRed = uint8(ConvertRange(0, 255, 0, 128, p.Red))
		r.LightBarGreen = uint8(ConvertRange(0, 255, 0, 128, p.Green))
		r.LightBarBlue = uint8(ConvertRange(0, 255, 0, 128, p.Blue))

	// flag1 bit 3 ?

	// flag1 bit 4
	case PlayerLEDs:
		//fmt.Printf("[Emit0x31] %#v\n", p)

		r.ValidFlag1 |= DS_OUTPUT_VALID_FLAG1_PLAYER_INDICATOR_CONTROL_ENABLE

		// LedID to Bit map
		led := []uint8{
			0: 1 << 0,
			1: 1 << 1,
			2: 1 << 2,
			3: 1 << 3,
			4: 1 << 4,
		}
		for id := 0; id < 5; id++ {
			if p[id] {
				r.PlayerLEDs |= led[id]
			}
		}

	// flag1 bit 4 - 7 ?

	// flag2 bit 0 ?

	// flag2 bit 1
	case LEDSetup:
		//fmt.Printf("[Emit0x31] %#v\n", p)

		r.ValidFlag2 = DS_OUTPUT_VALID_FLAG2_LIGHTBAR_SETUP_CONTROL_ENABLE
		r.LightBarSetup = DS_OUTPUT_LIGHTBAR_SETUP_LIGHT_OUT
	}

	// flag2 bit 2 - 7 ?
}
