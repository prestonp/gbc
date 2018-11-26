package cpu

import (
	"github.com/prestonp/gbc/mmu"
)

type Register uint8

const (
	A Register = iota
	F

	B
	C

	D
	E

	H
	L
)

type CPU struct {
	R []uint8

	SP uint16
	PC uint16

	M int // machine clock
	T int // instruction clock

	MMU mmu.MMU
}

func New() *CPU {
	return &CPU{
		SP: 0xfffe,
		PC: 0x100,

		R:   make([]uint8, 8),
		MMU: mmu.New(),
	}
}

func (c *CPU) Tick() {
	op := c.readByte()

	switch op {
	// 0X Range ///////////////////////////
	case 0x01:
		// LD BC, nn
		c.ldWordWord(B, C, c.readWord())
	case 0x02:
		// LD (BC), A
		c.ldAddrReg(toWord(c.R[B], c.R[C]), A)
	case 0x06:
		// LD B, n
		c.ldRegByte(B, c.readByte())
	case 0x08:
		// LD (nn), SP
		c.ldAddrWord(c.readWord(), c.SP)
	case 0x0A:
		// LD A, (BC)
		c.ldRegAddr(A, toWord(c.R[B], c.R[C]))
	case 0x0E:
		// LD C, n
		c.ldRegByte(C, c.readByte())

	// 1X Range ///////////////////////////
	case 0x11:
		// LD DE, nn
		c.ldWordWord(D, E, c.readWord())
	case 0x12:
		// LD (DE), A
		c.ldAddrReg(toWord(c.R[D], c.R[E]), A)
	case 0x16:
		// LD D, n
		c.ldRegByte(D, c.readByte())
	case 0x1A:
		// LD A, (DE)
		c.ldRegAddr(A, toWord(c.R[B], c.R[C]))
	case 0x1E:
		// LD E, n
		c.ldRegByte(E, c.readByte())

	// 2X Range ///////////////////////////
	case 0x21:
		// LD HL, nn
		c.ldWordWord(H, L, c.readWord())
	case 0x22:
		// LD (HLI), A
		// LD (HL+), A
		// LDI (HL), A
		c.ldAddrReg(toWord(c.R[H], c.R[L]), A).incWord(H, L)
	case 0x26:
		// LD H, n
		c.ldRegByte(H, c.readByte())
	case 0x2A:
		// LD A, (HLI)
		// LD A, (HL+)
		// LDI A, (HL)
		c.ldRegAddr(A, toWord(c.R[H], c.R[L])).incWord(H, L)
	case 0x2E:
		// LD L, n
		c.ldRegByte(L, c.readByte())

	// 3X Range ///////////////////////////
	case 0x31:
		// LD SP, nn
		c.ldSpWord(c.readWord())
	case 0x32:
		// LD (HLD), A
		// LD (HL-), A
		// LDD (HL), A
		c.
			ldAddrReg(toWord(c.R[H], c.R[L]), A).
			decWord(H, L)
	case 0x36:
		// LD (HL), n
		c.ldAddrByte(toWord(c.R[H], c.R[L]), c.readByte())
	case 0x3A:
		// LD A, (HLD)
		// LD A, (HL-)
		// LDD A, (HL)
		c.
			ldRegAddr(A, toWord(c.R[H], c.R[L])).
			decWord(H, L)
	case 0x3E:
		// LD A, n
		c.ldRegByte(A, c.readByte())
	// 4X Range ///////////////////////////
	case 0x40:
		// LD B, B
		c.ldRegReg(B, B)
	case 0x41:
		// LD B, C
		c.ldRegReg(B, C)
	case 0x42:
		// LD B, D
		c.ldRegReg(B, D)
	case 0x43:
		// LD B, E
		c.ldRegReg(B, E)
	case 0x44:
		// LD B, H
		c.ldRegReg(B, H)
	case 0x45:
		// LD B, L
		c.ldRegReg(B, L)
	case 0x46:
		// LD B, (HL)
		c.ldRegAddr(B, toWord(c.R[H], c.R[L]))
	case 0x47:
		// LD B, A
		c.ldRegReg(B, A)
	case 0x48:
		// LD C, B
		c.ldRegReg(C, B)
	case 0x49:
		// LD C, C
		c.ldRegReg(C, C)
	case 0x4A:
		// LD C, D
		c.ldRegReg(C, D)
	case 0x4B:
		// LD C, E
		c.ldRegReg(C, E)
	case 0x4C:
		// LD C, H
		c.ldRegReg(C, H)
	case 0x4D:
		// LD C, L
		c.ldRegReg(C, L)
	case 0x4E:
		// LD C, (HL)
		c.ldRegAddr(C, toWord(c.R[H], c.R[L]))
	case 0x4F:
		// LD C, A
		c.ldRegReg(C, A)

	// 5X Range ///////////////////////////
	case 0x50:
		// LD D, B
		c.ldRegReg(D, B)
	case 0x51:
		// LD D, C
		c.ldRegReg(D, C)
	case 0x52:
		// LD D, D
		c.ldRegReg(D, D)
	case 0x53:
		// LD D, E
		c.ldRegReg(D, E)
	case 0x54:
		// LD D, H
		c.ldRegReg(D, H)
	case 0x55:
		// LD D, L
		c.ldRegReg(D, L)
	case 0x56:
		// LD D, (HL)
		c.ldRegAddr(D, toWord(c.R[H], c.R[L]))
	case 0x57:
		// LD D, A
		c.ldRegReg(D, A)
	case 0x58:
		// LD E, B
		c.ldRegReg(E, B)
	case 0x59:
		// LD E, C
		c.ldRegReg(E, C)
	case 0x5A:
		// LD E, D
		c.ldRegReg(E, D)
	case 0x5B:
		// LD E, E
		c.ldRegReg(E, E)
	case 0x5C:
		// LD E, H
		c.ldRegReg(E, H)
	case 0x5D:
		// LD E, L
		c.ldRegReg(E, L)
	case 0x5E:
		// LD E, (HL)
		c.ldRegAddr(E, toWord(c.R[H], c.R[L]))
	case 0x5F:
		// LD E, A
		c.ldRegReg(E, A)

	// 6X Range ///////////////////////////
	case 0x60:
		// LD H, B
		c.ldRegReg(H, B)
	case 0x61:
		// LD H, C
		c.ldRegReg(H, C)
	case 0x62:
		// LD H, D
		c.ldRegReg(H, D)
	case 0x63:
		// LD H, E
		c.ldRegReg(H, E)
	case 0x64:
		// LD H, H
		c.ldRegReg(H, H)
	case 0x65:
		// LD H, L
		c.ldRegReg(H, L)
	case 0x66:
		// LD H, (HL)
		c.ldRegAddr(H, toWord(c.R[H], c.R[L]))
	case 0x67:
		// LD H, A
		c.ldRegReg(H, A)
	case 0x68:
		// LD L, B
		c.ldRegReg(L, B)
	case 0x69:
		// LD L, C
		c.ldRegReg(L, C)
	case 0x6A:
		// LD L, D
		c.ldRegReg(L, D)
	case 0x6B:
		// LD L, E
		c.ldRegReg(L, E)
	case 0x6C:
		// LD L, H
		c.ldRegReg(L, H)
	case 0x6D:
		// LD L, L
		c.ldRegReg(L, L)
	case 0x6E:
		// LD L, (HL)
		c.ldRegAddr(L, toWord(c.R[H], c.R[L]))
	case 0x6F:
		// LD L, A
		c.ldRegReg(L, A)

	// 7X Range ///////////////////////////
	case 0x70:
		// LD (HL), B
		c.ldAddrReg(toWord(c.R[H], c.R[L]), B)
	case 0x71:
		// LD (HL), C
		c.ldAddrReg(toWord(c.R[H], c.R[L]), C)
	case 0x72:
		// LD (HL), D
		c.ldAddrReg(toWord(c.R[H], c.R[L]), D)
	case 0x73:
		// LD (HL), E
		c.ldAddrReg(toWord(c.R[H], c.R[L]), E)
	case 0x74:
		// LD (HL), H
		c.ldAddrReg(toWord(c.R[H], c.R[L]), H)
	case 0x75:
		// LD (HL), L
		c.ldAddrReg(toWord(c.R[H], c.R[L]), L)
	case 0x77:
		// LD (HL), A
		c.ldAddrReg(toWord(c.R[H], c.R[L]), A)
	case 0x7F:
		// LD A, A
		c.ldRegReg(A, A)
	case 0x78:
		// LD A, B
		c.ldRegReg(A, B)
	case 0x79:
		// LD A, C
		c.ldRegReg(A, C)
	case 0x7A:
		// LD A, D
		c.ldRegReg(A, D)
	case 0x7B:
		// LD A, E
		c.ldRegReg(A, E)
	case 0x7C:
		// LD A, H
		c.ldRegReg(A, H)
	case 0x7D:
		// LD A, L
		c.ldRegReg(A, L)
	case 0x7E:
		// LD A, (HL)
		c.ldRegAddr(A, toWord(c.R[H], c.R[L]))

	// EX Range ///////////////////////////
	case 0xE0:
		// LDH (n), A
		c.ldAddrReg(toWord(0xFF, c.readByte()), A)
	case 0xE2:
		// LD (C), A
		c.ldAddrReg(toWord(0xFF, c.R[C]), A)
	case 0xEA:
		// LD (nn), A
		c.ldAddrReg(c.readWord(), A)

	// FX Range ///////////////////////////
	case 0xF0:
		// LDH A, (n)
		c.ldRegAddr(A, toWord(0xFF, c.readByte()))
	case 0xF2:
		// LD A, (C)
		c.ldRegAddr(A, toWord(0xFF, c.R[C]))
	case 0xF9:
		// LD SP, HL
		// Added artifical internal delay: https://github.com/Gekkio/mooneye-gb/blob/master/docs/accuracy.markdown#some-instructions-take-more-cycles-than-just-the-memory-accesses-at-which-point-in-the-instruction-execution-do-these-extra-cycles-occur
		c.ldSpWord(toWord(c.R[H], c.R[L])).delay(1, 4)
	case 0xFA:
		// LD A, (nn)
		c.ldRegAddr(A, c.readWord())
	case 0xF8:
		// LD HL, SP+n
		// LDHL SP, n
		// Added artifical internal delay: https://github.com/Gekkio/mooneye-gb/blob/master/docs/accuracy.markdown#some-instructions-take-more-cycles-than-just-the-memory-accesses-at-which-point-in-the-instruction-execution-do-these-extra-cycles-occur
		c.ldHLSPn().delay(1, 4)
	}
}

// Load Register variants
func (c *CPU) ldRegByte(dst Register, b uint8) *CPU {
	c.R[dst] = b
	return c
}

func (c *CPU) ldRegReg(dst, src Register) *CPU {
	c.R[dst] = c.R[src]
	return c
}

func (c *CPU) ldRegAddr(dst Register, addr uint16) *CPU {
	c.R[dst] = c.MMU.ReadByte(addr)
	c.M++
	c.T += 4
	return c
}

// Load word variants
func (c *CPU) ldWordWord(hi, lo Register, word uint16) *CPU {
	c.R[hi] = uint8(word >> 8)
	c.R[lo] = uint8(word & 0xFF)
	return c
}

// Load address variants
func (c *CPU) ldAddrReg(addr uint16, src Register) *CPU {
	c.MMU.WriteByte(addr, c.R[src])
	c.M++
	c.T += 4
	return c
}

func (c *CPU) ldAddrByte(addr uint16, b uint8) *CPU {
	c.MMU.WriteByte(addr, b)
	c.M++
	c.T += 4
	return c
}

func (c *CPU) ldAddrWord(addr uint16, w uint16) *CPU {
	c.MMU.WriteByte(addr, uint8(w&0xFF))
	c.MMU.WriteByte(addr+1, uint8(w&0xFF00>>8))
	c.M += 2
	c.T += 8
	return c
}

// Load stack pointer variants
func (c *CPU) ldSpWord(word uint16) *CPU {
	c.SP = word
	return c
}

func (c *CPU) ldHLSPn() *CPU {
	// fetch signed byte from PC
	n := uint16(int16(c.readByte()))
	sp := c.SP + n

	// reset F register
	c.R[F] = 0

	// detect half carry
	if ((c.SP&0xF)+(n&0xF))&0x10 == 0x10 {
		c.R[F] |= 0x20
	}

	// detect carry
	if ((c.SP&0xFF)+(n&0xFF))&0x100 == 0x100 {
		c.R[F] |= 0x10
	}

	return c.ldWordWord(H, L, sp)
}

// read a byte from the PC a.k.a `n`
func (c *CPU) readByte() uint8 {
	b := c.MMU.ReadByte(c.PC)
	c.PC++
	c.M++
	c.T += 4
	return b
}

// read a word from the PC a.k.a `nn`
func (c *CPU) readWord() uint16 {
	return toWord(c.readByte(), c.readByte())
}

func (c *CPU) decWord(upper, lower Register) *CPU {
	w := toWord(c.R[upper], c.R[lower]) - 1
	c.R[upper] = uint8(w >> 8)
	c.R[lower] = uint8(w & 0xFF)
	return c
}

func (c *CPU) incWord(upper, lower Register) *CPU {
	w := toWord(c.R[upper], c.R[lower]) + 1
	c.R[upper] = uint8(w >> 8)
	c.R[lower] = uint8(w & 0xFF)
	return c
}

func (c *CPU) delay(m, t int) *CPU {
	c.M += m
	c.T += t
	return c
}

// combine two bytes into a word such that a is upper
// half and b is lower half
func toWord(a, b uint8) uint16 {
	return uint16(a)<<8 | uint16(b)
}
