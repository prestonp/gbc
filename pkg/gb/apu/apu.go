package apu

import "log"

type APU struct {
	nr52 byte
	nr11 byte
}

func New() *APU {
	return &APU{}
}

func (a *APU) SetRegister(addr uint16, b byte) {
	switch {
	case addr == 0xFF11:
		a.nr11 = b
	case addr == 0xFF26:
		a.nr52 = b & 0x8F
	default:
		log.Fatalf("unimplemented sound controller register 0x%04X = 0x%02X\n", addr, b)
	}
}

func (a *APU) GetRegister(addr uint16) byte {
	switch {
	case addr == 0xFF11:
		return a.nr11
	case addr == 0xFF26:
		return a.nr52
	default:
		log.Fatalf("unimplemented sound controller register 0x%04X\n", addr)
	}
	panic("unexpected sound controller failure")
}
