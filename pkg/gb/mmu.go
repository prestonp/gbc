package gb

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type ByteFlag byte

func (bf ByteFlag) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "\tVBlank: %t\n", bf&BitVBlank > 0)
	fmt.Fprintf(&b, "\tLCDStat: %t\n", bf&BitLCDStat > 0)
	fmt.Fprintf(&b, "\tTimer: %t\n", bf&BitTimer > 0)
	fmt.Fprintf(&b, "\tSerial: %t\n", bf&BitSerial > 0)
	fmt.Fprintf(&b, "\tJoypad: %t\n", bf&BitJoypad > 0)
	return b.String()
}

type MMU struct {
	booted bool // $00-$FF point to cartridge after booting
	boot   []byte
	rom    []byte
	wram   []byte
	hram   []byte
	IF     ByteFlag
	IE     ByteFlag
	SB     byte
	SC     byte
	BGP    byte // background and window palette

	gpu Module
	apu Module
}

func NewMMU(bootRom, cartRom []uint8, gpu, apu Module) *MMU {
	return &MMU{
		boot: bootRom,
		rom:  cartRom,
		wram: make([]byte, 8*1024),
		hram: make([]byte, 256),
		IF:   0,
		gpu:  gpu,
		apu:  apu,
	}
}

func ReadRom(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (m *MMU) ReadByte(a uint16) byte {
	switch {
	case a >= 0x0000 && a < 0x8000:
		if a <= 0xFF && !m.booted {
			return m.boot[a]
		}
		return m.rom[a]
	case a >= 0x8000 && a < 0xA000:
		m.gpu.ReadByte(a)
	case a >= 0xC000 && a < 0xE000:
		// working ram
		return m.wram[a-0xC000]
	case a >= 0xE000 && a < 0xFE00:
		// echo ram
		return m.wram[a-0xE000]
	case a >= 0xFE00 && a <= 0xFE9F:
		m.gpu.ReadByte(a)
	case a >= 0xFEA0 && a <= 0xFEFF:
		return 0x00
	case a == 0xFF01:
		// SB - serial transfer data
		return m.SB
	case a == 0xFF02:
		// SC - serial transfer control
		return m.SC
	case a == 0xFF0F:
		// IF Interrupt flag
		return byte(m.IF)
	case a >= 0xFF10 && a <= 0xFF26:
		return m.apu.ReadByte(a)
	case a >= 0xFF40 && a <= 0xFF48:
		return m.gpu.ReadByte(a)
	case a >= 0xFF80 && a < 0xFFFF:
		return m.hram[a-0xFF80]
	case a == 0xFFFF:
		// IE Interrupt enable
		return byte(m.IE)
	default:
		log.Panicf("unimplemented mmu address: 0x%04X", a)
	}
	return 0
}

func (m *MMU) WriteByte(a uint16, n uint8) {
	switch {
	case a >= 0x0000 && a <= 0x3FFF:
		// this is rom but tetris will try to
		// write to it, skip this
		fmt.Printf("warning: skipping over write to rom 0x%04X 0x%02X\n", a, n)
	case a >= 0x8000 && a < 0xA000:
		m.gpu.WriteByte(a, n)
	case a >= 0xC000 && a <= 0xE000:
		// working ram
		m.wram[a-0xC000] = n
	case a >= 0xE000 && a <= 0xFE00:
		// echo of working ram
		m.wram[a-0xE000] = n
	case a == 0xFF01:
		// SB - serial transfer data
		m.SB = n
	case a == 0xFF02:
		// SC - serial transfer control
		// todo: not implemented
		m.SC = n
	case a == 0xFF0F:
		// IF - Interrupt Flag
		m.IF = ByteFlag(n)
	case a >= 0xFF10 && a <= 0xFF26:
		m.apu.WriteByte(a, n)
	case a >= 0xFF40 && a <= 0xFF49:
		m.gpu.WriteByte(a, n)
	case a == 0xFF50:
		m.booted = n != 0
	case a >= 0xFF80 && a < 0xFFFF:
		// high ram
		m.hram[a-0xFF80] = n
	case a >= 0xFE00 && a <= 0xFE9F:
		m.gpu.WriteByte(a, n)
	case a >= 0xFEA0 && a <= 0xFEFF:
		// unusable area
	case a >= 0xFF4C && a <= 0xFF7F:
		// unusable area
	case a == 0xFFFF:
		// IE - Interrupt Enable
		m.IE = ByteFlag(n)
	default:
		log.Panicf("cannot write to mmu 0x%04X = 0x%02X", a, n)
	}
}

const (
	BitVBlank ByteFlag = 1 << iota
	BitLCDStat
	BitTimer
	BitSerial
	BitJoypad
)
