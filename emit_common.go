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

type LEDSetup struct{}

type OutputCommon struct {
	// bit 0: [Rumble] DS_OUTPUT_VALID_FLAG0_COMPATIBLE_VIBRATION
	// bit 1: [Rumble] DS_OUTPUT_VALID_FLAG0_HAPTICS_SELECT
	// bit 2-7: TODO
	ValidFlag0 byte

	// bit 0: [Mic] DS_OUTPUT_VALID_FLAG1_MIC_MUTE_LED_CONTROL_ENABLE
	// bit 1: [Mic] DS_OUTPUT_VALID_FLAG1_POWER_SAVE_CONTROL_ENABLE
	// bit 2: [LightBar] DS_OUTPUT_VALID_FLAG1_LIGHTBAR_CONTROL_ENABLE
	// bit 3: TODO
	// bit 4: [PlayerLEDs] DS_OUTPUT_VALID_FLAG1_PLAYER_INDICATOR_CONTROL_ENABLE
	// bit 5-7: TODO
	ValidFlag1 byte

	// [Rumble] uint8
	MotorRight byte
	// [Rumble] uint8
	MotorLeft byte

	// TODO
	Reserved [4]byte

	// bit 0: Mute LED
	// bit 1-7: TODO
	MuteButtonLED byte

	// bit 0-3: TODO
	// bit 4: [Mic] DS_OUTPUT_POWER_SAVE_CONTROL_MIC_MUTE
	// bit 5-7: TODO
	PowerSaveControl byte

	// TODO
	Reserved2 [28]byte

	// bit 0: TODO
	// bit 1: [LEDSetup] DS_OUTPUT_VALID_FLAG2_LIGHTBAR_SETUP_CONTROL_ENABLE
	// bit 2-7: TODO
	ValidFlag2 byte

	// TODO
	Reserved3 [2]byte

	// bit 0: TODO
	// bit 1: [LEDSetup] DS_OUTPUT_LIGHTBAR_SETUP_LIGHT_OUT
	// bit 2-7: TODO
	LightBarSetup byte

	// TODO
	LedBrightness byte

	// bit 0-4: [PlayerLEDs] 5 led's
	// bit 5: [PlayerLEDs] fade LED animation 0:on/1:off (default: 0)
	// bit 6-7: TODO
	PlayerLEDs byte

	// [LightBar] uint8
	LightBarRed byte
	// [LightBar] uint8
	LightBarGreen byte
	// [LightBar] uint8
	LightBarBlue byte
}

func (r *OutputCommon) ApplyProp(prop interface{}) {
	switch prop := prop.(type) {
	case Rumble:
		//fmt.Printf("[Emit0x31] %#v\n", p)

		r.ValidFlag0 |= DS_OUTPUT_VALID_FLAG0_COMPATIBLE_VIBRATION
		r.ValidFlag0 |= DS_OUTPUT_VALID_FLAG0_HAPTICS_SELECT

		r.MotorRight = prop.Right
		r.MotorLeft = prop.Left
	case Mic:
		//fmt.Printf("[Emit0x31] %#v\n", p)

		// enable/disable mute LED
		r.ValidFlag1 |= DS_OUTPUT_VALID_FLAG1_MIC_MUTE_LED_CONTROL_ENABLE
		if prop.LED {
			r.MuteButtonLED = 0x1
		} else {
			r.MuteButtonLED = 0x0
		}

		// enable/disable microphone
		r.ValidFlag1 |= DS_OUTPUT_VALID_FLAG1_POWER_SAVE_CONTROL_ENABLE
		if prop.Muted {
			r.PowerSaveControl |= DS_OUTPUT_POWER_SAVE_CONTROL_MIC_MUTE
		} else {
			r.PowerSaveControl &= ^DS_OUTPUT_POWER_SAVE_CONTROL_MIC_MUTE
		}
	case LightBar:
		//fmt.Printf("[Emit0x31] %#v\n", p)

		r.ValidFlag1 |= DS_OUTPUT_VALID_FLAG1_LIGHTBAR_CONTROL_ENABLE

		// convert 0-255 to 0-128
		r.LightBarRed = uint8(ConvertRange(0, 255, 0, 128, prop.Red))
		r.LightBarGreen = uint8(ConvertRange(0, 255, 0, 128, prop.Green))
		r.LightBarBlue = uint8(ConvertRange(0, 255, 0, 128, prop.Blue))
	case PlayerLEDs:
		//fmt.Printf("[Emit0x31] %#v\n", p)

		r.ValidFlag1 |= DS_OUTPUT_VALID_FLAG1_PLAYER_INDICATOR_CONTROL_ENABLE

		// led bit map
		ledBitMap := []uint8{
			0: 1 << 0,
			1: 1 << 1,
			2: 1 << 2,
			3: 1 << 3,
			4: 1 << 4,
		}
		for ledID := 0; ledID < 5; ledID++ {
			if prop.LED[ledID] {
				r.PlayerLEDs |= ledBitMap[ledID]
			}
		}

		// fade-animate led change
		if prop.DisableChangeAnimation {
			r.PlayerLEDs |= 1 << 5
		} else {
			r.PlayerLEDs |= 0 << 5
		}
	case LEDSetup:
		//fmt.Printf("[Emit0x31] %#v\n", p)

		r.ValidFlag2 = DS_OUTPUT_VALID_FLAG2_LIGHTBAR_SETUP_CONTROL_ENABLE
		r.LightBarSetup = DS_OUTPUT_LIGHTBAR_SETUP_LIGHT_OUT
	}
}
