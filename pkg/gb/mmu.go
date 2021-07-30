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
	boot []byte
	rom  []byte
	wram []byte
	vram []byte
	hram []byte
	IF   ByteFlag
	IE   ByteFlag
	SB   byte
	SC   byte
	BGP  byte // background and window palette

	gpu GPU
	apu Module
}

func NewMMU(bootRom, cartRom []uint8, gpu GPU, apu Module) *MMU {
	return &MMU{
		boot: bootRom,
		rom:  cartRom,
		wram: make([]byte, 8*1024),
		vram: make([]byte, 8*1024),
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
		if a <= 0xFF {
			return m.boot[a]
		}
		return m.rom[a]
	case a >= 0x8000 && a < 0xA000:
		return m.vram[a-0x8000]
	case a >= 0xC000 && a < 0xE000:
		// working ram
		return m.wram[a-0xC000]
	case a >= 0xE000 && a < 0xFE00:
		// echo ram
		return m.wram[a-0xE000]
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
		return m.apu.GetRegister(a)
	case a == 0xFF40:
		// LCD Control
		return m.gpu.GetControl()
	case a == 0xFF41:
		return m.gpu.GetStat()
	case a == 0xFF42:
		return m.gpu.GetScrollY()
	case a == 0xFF43:
		return m.gpu.GetScrollX()
	case a == 0xFF44:
		return m.gpu.GetLY()
	case a == 0xFF47:
		return m.gpu.GetRegister(a)
	case a >= 0xFF80 && a < 0xFFFF:
		return m.hram[a-0xFF80]
	case a == 0xFFFF:
		// IE Interrupt enable
		return byte(m.IE)
	default:
		log.Panicf("unimplemented memory address: 0x%04X", a)
	}
	return 0
}

func (m *MMU) WriteByte(a uint16, n uint8) {
	switch {
	case a >= 0x8000 && a < 0xA000:
		// video ram
		m.vram[a-0x8000] = n
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
		m.apu.SetRegister(a, n)
	case a == 0xFF40:
		// LCD Control
		m.gpu.SetControl(n)
	case a == 0xFF41:
		m.gpu.SetStat(n & 0xF8)
	case a == 0xFF42:
		m.gpu.SetScrollY(n)
	case a == 0xFF43:
		m.gpu.SetScrollX(n)
	case a == 0xFF44:
		m.gpu.ResetLY()
	case a == 0xFF47:
		m.gpu.SetRegister(a, n)
	case a >= 0xFF80 && a < 0xFFFF:
		// high ram
		m.hram[a-0xFF80] = n
	case a == 0xFFFF:
		// IE - Interrupt Enable
		m.IE = ByteFlag(n)
	default:
		log.Panicf("cannot write 0x%04X = 0x%02X", a, n)
	}
}

const (
	BitVBlank ByteFlag = 1 << iota
	BitLCDStat
	BitTimer
	BitSerial
	BitJoypad
)
