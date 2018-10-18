package A75C4269

const (
	PowerOff byte = iota
	PowerOn
	PowerOnAndOffTimer
	PowerOffAndOnTimer

	ModeCooler byte = iota
	ModeHeater
	ModeDehumidifier

	WindDirectionAuto byte = iota
	WindDirection1
	WindDirection2
	WindDirection3
	WindDirection4
	WindDirection5

	AirVolumeAuto byte = iota
	AirVolumeStill
	AirVolume1
	AirVolume2
	AirVolume3
	AirVolume4
	AirVolumePowerful
)

type Controller struct {
	Power         byte
	Mode          byte
	PresetTemp    uint
	AirVolume     byte
	WindDirection byte
	TimerHour     byte
}

func (r *Controller) GetSignalBytes() (result []byte) {
	var b byte

	// 1-5byte: header
	result = append(result, 0x02, 0x20, 0x0E, 0x04, 0x00)

	// 6byte: mode, timer, power
	b = 0
	switch r.Mode {
	case ModeCooler:
		b |= 0x3 // 0011
	case ModeHeater:
		b |= 0x4 // 0100
	case ModeDehumidifier:
		b |= 0x2 // 0010
	default:
		b |= 0x3 // Cooler
	}
	b <<= 4
	switch r.Power {
	case PowerOff:
		b |= 0x0 // 0000
	case PowerOn:
		b |= 0x1 // 0001
	case PowerOnAndOffTimer:
		b |= 0x5 // 0101
	case PowerOffAndOnTimer:
		b |= 0x2 // 0010
	default:
		b |= 0x0 // PowerOff
	}
	result = append(result, b)

	// 7byte: temp
	b = 0x20 // 00100000
	switch {
	case r.PresetTemp < 16:
		b |= 0x0 << 1 // 0000
	case r.PresetTemp > 30:
		b |= 0xE << 1 // 1110
	default:
		b |= byte(r.PresetTemp-16) << 1
	}
	result = append(result, b)

	// 8byte: ?
	b = 0x80
	result = append(result, b)

	// 9byte: air volume, wind direction
	switch r.AirVolume {
	case AirVolumeAuto:
		b = 0xA << 4 // 1010
	case AirVolume1, AirVolumeStill, AirVolumePowerful:
		b = 0x3 << 4 // 0011
	case AirVolume2:
		b = 0x4 << 4 // 0101
	case AirVolume3:
		b = 0x5 << 4 // 0110
	case AirVolume4:
		b = 0x6 << 4 // 0111
	default:
		b = 0xA << 4 // Auto
	}
	switch r.WindDirection {
	case WindDirectionAuto:
		b |= 0xF // 1111
	case WindDirection1:
		b |= 0x1 // 0001
	case WindDirection2:
		b |= 0x2 // 0010
	case WindDirection3:
		b |= 0x3 // 0011
	case WindDirection4:
		b |= 0x4 // 0100
	case WindDirection5:
		b |= 0x5 // 0101
	default:
		b |= 0xF // Auto
	}
	result = append(result, b)

	// 10byte: ?
	b = 0x00
	result = append(result, b)

	// 11byte: timer setting
	b = 0
	if r.Power == PowerOffAndOnTimer || r.Power == PowerOnAndOffTimer {
		b = 0x3C
	}
	result = append(result, b)

	// 12-13byte: timer hour
	if r.Power == PowerOnAndOffTimer || r.Power == PowerOffAndOnTimer {
		switch r.TimerHour {
		case 1:
			result = append(result, 0xC0, 0x03) // 11000000 00000011
		case 2:
			result = append(result, 0x80, 0x07) // 10000000 00000111
		case 3:
			result = append(result, 0x40, 0x0B) // 01000000 00001011
		case 4:
			result = append(result, 0x00, 0x0F) // 00000000 00001111
		case 5:
			result = append(result, 0xC0, 0x12) // 11000000 00010010
		case 6:
			result = append(result, 0x80, 0x16) // 10000000 00010110
		case 7:
			result = append(result, 0x40, 0x1A) // 01000000 00011010
		case 8:
			result = append(result, 0x00, 0x1E) // 00000000 00011110
		case 9:
			result = append(result, 0xC0, 0x21) // 11000000 00100001
		case 10:
			result = append(result, 0x80, 0x25) // 10000000 00100101
		case 11:
			result = append(result, 0x40, 0x29) // 01000000 00101001
		case 12:
			result = append(result, 0x00, 0x2D) // 00000000 00101101
		default:
			result = append(result, 0xC0, 0x03) // 1 hour
		}
	} else {
		result = append(result, 0x06, 0x60) // 00000110 01100000
	}

	// 14byte: Various Flags
	b = 0
	if r.AirVolume == AirVolumeStill {
		b |= 1
	}
	b <<= 5
	if r.AirVolume == AirVolumePowerful {
		b |= 1
	}
	result = append(result, b)

	// 15byte: Various Flags
	b = 0
	if r.PresetTemp <= 16 || r.PresetTemp >= 30 {
		b |= 1
	}
	b <<= 1
	result = append(result, b)

	// 16byte: ?
	b = 0x80
	result = append(result, b)

	// 17byte: ?
	b = 0x00
	result = append(result, b)

	// 18byte: ?
	b = 0x06
	result = append(result, b)

	// 19byte: checksum
	sum := 0x6
	for i := 5; i < 18; i++ {
		sum += int(result[i])
	}
	b = byte(0xFF & sum)
	result = append(result, b)

	return result
}
