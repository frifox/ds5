package ds5

type PlayerLEDs [5]bool

func (l *PlayerLEDs) AllOff() {
	for id, _ := range l {
		l[id] = false
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
