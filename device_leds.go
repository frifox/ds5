package ds5

type PlayerLEDs [5]bool

func (l *PlayerLEDs) AllOff() {
	for id, _ := range l {
		l[id] = false
	}
}

//  mimic PS5 behavior
var ledPlayerMap = map[uint8][]uint8{
	0: {},
	1: {2},
	2: {1, 3},
	3: {0, 2, 4},
	4: {0, 1, 3, 4},
	5: {0, 1, 2, 3, 4},
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
	0: {},
	1: {0},
	2: {0, 1},
	3: {0, 1, 2},
	4: {0, 1, 2, 3},
	5: {0, 1, 2, 3, 4},
}

func (l *PlayerLEDs) SetBar(count uint8) {
	l.AllOff()
	if leds, exists := ledBarMap[count]; exists {
		for _, id := range leds {
			l[id] = true
		}
	}
}
