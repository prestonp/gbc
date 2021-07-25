package apu

import "log"

type APU struct {
	nr10 byte
	nr11 byte
	nr12 byte
	nr50 byte // unimplemented
	nr51 byte // unimplemented
	nr52 byte
}

func New() *APU {
	return &APU{}
}

func (a *APU) SetRegister(addr uint16, b byte) {
	switch {
	case addr == 0xFF10:
		a.nr10 = b & 0x7F
	case addr == 0xFF11:
		a.nr11 = b
	case addr == 0xFF12:
		a.nr12 = b
	case addr == 0xFF24:
		a.nr50 = b
	case addr == 0xFF25:
		a.nr51 = b
	case addr == 0xFF26:
		a.nr52 = b & 0x8F
	default:
		log.Fatalf("unimplemented sound controller register 0x%04X = 0x%02X\n", addr, b)
	}
}

func (a *APU) GetRegister(addr uint16) byte {
	switch {
	case addr == 0xFF10:
		return a.nr10
	case addr == 0xFF11:
		return a.nr11
	case addr == 0xFF12:
		return a.nr12
	case addr == 0xFF24:
		return a.nr50
	case addr == 0xFF25:
		return a.nr51
	case addr == 0xFF26:
		return a.nr52
	default:
		log.Fatalf("unimplemented sound controller register 0x%04X\n", addr)
	}
	panic("unexpected sound controller failure")
}

// 000: sweep off - no freq change 001: 7.8 ms (1/128Hz)
// 010: 15.6 ms (2/128Hz)
// 011: 23.4 ms (3/128Hz)
func (a *APU) sweepTime() byte {
	return a.nr10 >> 4
}

// false: Addition (frequency increases)
// true: Subtraction (frequency decreases)
func (a *APU) sweepMode() bool {
	return a.nr10&(1<<3) > 0
}

func (a *APU) sweepShift() byte {
	return a.nr10 & 0x7
}

func (a *APU) envelopeInitVolume() byte {
	return a.nr12 >> 4
}

// false: attenuate
// true: amplify
func (a *APU) envelopeMode() bool {
	return a.nr12&0x8 == 0x8
}

// false: attenuate
// true: amplify
func (a *APU) envelopeSweep() byte {
	return a.nr12 & 0x7
}
