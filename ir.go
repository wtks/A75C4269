package A75C4269

const (
	T           = 445
	T8          = T * 8
	T4          = T * 4
	T3          = T * 3
	TracerSpace = T * 20
)

var (
	callSign = []byte{0x02, 0x20, 0x0E, 0x04, 0x00, 0x00, 0x00, 0x06}
)

func (r *Controller) GetRawSignal() []uint32 {
	var d []uint32
	d = append(d, appendTracerSpace(convertRawSignal(callSign))...)
	d = append(d, convertRawSignal(r.GetSignalBytes())...)
	return d
}

func convertRawSignal(bytes []byte) []uint32 {
	var seq []uint32

	// Leader
	seq = append(seq, T8, T4)

	// Customer Code 1, 2
	seq = appendByte(seq, bytes[0])
	seq = appendByte(seq, bytes[1])

	// Parity & Data0
	var i uint
	for i = 4; i < 8; i++ {
		if refBit(bytes[2], i) == 1 {
			seq = appendBit1(seq)
		} else {
			seq = appendBit0(seq)
		}
	}
	for i = 0; i < 4; i++ {
		if refBit(bytes[2], i) == 1 {
			seq = appendBit1(seq)
		} else {
			seq = appendBit0(seq)
		}
	}

	// DataN
	for i := 3; i < len(bytes); i++ {
		seq = appendByte(seq, bytes[i])
	}

	// Tracer
	seq = append(seq, T)

	return seq
}

func appendTracerSpace(seq []uint32) []uint32 {
	return append(seq, TracerSpace)
}

func appendBit0(seq []uint32) []uint32 {
	return append(seq, T, T)
}

func appendBit1(seq []uint32) []uint32 {
	return append(seq, T, T3)
}

func refBit(i byte, b uint) byte {
	return (i >> b) & 1
}

func appendByte(seq []uint32, b byte) []uint32 {
	var i uint
	for i = 0; i < 8; i++ {
		if refBit(b, i) == 1 {
			seq = appendBit1(seq)
		} else {
			seq = appendBit0(seq)
		}
	}
	return seq
}
