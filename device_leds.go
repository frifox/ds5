package ds5

type PlayerLEDs [5]bool

func (l *PlayerLEDs) AllOff() {
	for id, _ := range l {
		l[id] = false
	}
}
func (l *PlayerLEDs) AllOn() {
	for id, _ := range l {
		l[id] = true
	}
}

//  mimic PS5 behavior
var ledPlayerMap = map[uint8][]uint8{
	0: {},              // _____
	1: {2},             // __X__
	2: {1, 3},          // _X_X_
	3: {0, 2, 4},       // X_X_X
	4: {0, 1, 3, 4},    // XX_XX
	5: {0, 1, 2, 3, 4}, // XXXXX
}

func (l *PlayerLEDs) SetPlayer(playerID uint8) {
	l.AllOff()
	if leds, exists := ledPlayerMap[playerID]; exists {
		for _, id := range leds {
			l[id] = true
		}
	}
}

// volume bar
var ledBarMap = map[uint8][]uint8{
	0: {},              // _____
	1: {0},             // X____
	2: {0, 1},          // XX___
	3: {0, 1, 2},       // XXX__
	4: {0, 1, 2, 3},    // XXXX_
	5: {0, 1, 2, 3, 4}, // XXXXX
}

func (l *PlayerLEDs) SetBar(count uint8) {
	l.AllOff()
	if leds, exists := ledBarMap[count]; exists {
		for _, id := range leds {
			l[id] = true
		}
	}
}

// marquee
var dotMap = map[uint8][]uint8{
	0: {},  // _____
	1: {0}, // X____
	2: {1}, // _X___
	3: {2}, // __X__
	4: {3}, // ___X_
	5: {4}, // ____X
}

func (l *PlayerLEDs) SetDot(count uint8) {
	l.AllOff()
	if leds, exists := dotMap[count]; exists {
		for _, id := range leds {
			l[id] = true
		}
	}
}
