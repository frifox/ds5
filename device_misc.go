package ds5

type OutputSequencer uint8

func (s *OutputSequencer) Get() uint8 {
	seq := *s // return current seq

	*s++          // next seq
	*s = *s & 0xf // truncate it to uint4 (0 to 15)

	return uint8(seq)
}
