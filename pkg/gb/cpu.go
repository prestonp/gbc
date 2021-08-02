package gb

import (
	"fmt"
	"log"
	"strings"

	"github.com/prestonp/gbc/pkg/logbuf"
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

func (r Register) String() string {
	switch r {
	case A:
		return "A"
	case F:
		return "F"

	case B:
		return "B"
	case C:
		return "C"

	case D:
		return "D"
	case E:
		return "E"

	case H:
		return "H"
	case L:
		return "L"
	default:
		log.Panicf("unknown register: %d", r)
	}
	return ""
}

type CPU struct {
	R []uint8

	SP uint16
	PC uint16

	M int // machine clock
	T int // instruction clock

	MMU *MMU
	GPU GPU

	debug bool

	IME      bool // interrupt master enable
	shouldDI bool // disable interrupts
	shouldEI bool // enable interrupts

	log *logbuf.Buffer
}

var (
	FlagZero      byte = 0x80
	FlagSubtract  byte = 0x40
	FlagHalfCarry byte = 0x20
	FlagCarry     byte = 0x10
)

func (c *CPU) Debugf(s string, args ...interface{}) {
	if !c.debug {
		return
	}

	fmt.Fprintf(c.log, "[debug] "+s, args...)
}

func NewCPU(mmu *MMU, gpu GPU, debug bool) *CPU {
	return &CPU{
		SP: 0x0,
		PC: 0x0,

		R:     make([]uint8, 8),
		MMU:   mmu,
		GPU:   gpu,
		debug: debug,

		log: logbuf.New(1024),
	}
}

func (c *CPU) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "State\n")
	fmt.Fprintf(&b, "SP:\t0x%04X\n", c.SP)
	fmt.Fprintf(&b, "PC:\t0x%04X\n", c.PC)
	fmt.Fprintf(&b, "A:\t0x%02X\n", c.R[A])
	fmt.Fprintf(&b, "F:\t0x%02X (ZNHC = %04b)\n", c.R[F], c.R[F]>>4)
	fmt.Fprintf(&b, "B:\t0x%02X\n", c.R[B])
	fmt.Fprintf(&b, "C:\t0x%02X\n", c.R[C])
	fmt.Fprintf(&b, "D:\t0x%02X\n", c.R[D])
	fmt.Fprintf(&b, "E:\t0x%02X\n", c.R[E])
	fmt.Fprintf(&b, "H:\t0x%02X\n", c.R[H])
	fmt.Fprintf(&b, "L:\t0x%02X\n", c.R[L])
	fmt.Fprintf(&b, "IME:\t%v\n", c.IME)
	fmt.Fprintf(&b, "IE:\n%s", c.MMU.IE)
	fmt.Fprintf(&b, "IF:\n%s", c.MMU.IF)
	fmt.Fprintf(&b, "PPU:\n%s\n", c.MMU.gpu)
	return b.String()
}

func (c *CPU) Run() {
	c.GPU.Loop(c.Update)
}

func (c *CPU) Update() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(c.log.String())
			fmt.Println(r)
		}
	}()

	c.resolveInterruptToggle()
	op := c.fetch()
	exec := c.decode(op)
	exec(c)
	c.Debugf("%s\n", c)

	if c.PC == 0x2817 {
		panic("STOP LOADING TILE DATA")
	}
}

// See DI/EI opcode reference for more context, but basically the effects of EI/DI instructions are delayed by
// one instruction
func (c *CPU) resolveInterruptToggle() {
	if c.shouldDI {
		c.IME = false
		c.shouldDI = false
	} else if c.shouldEI {
		c.IME = true
		c.shouldEI = false
	}
}

func (c *CPU) fetch() byte {
	op := c.readByte()
	c.Debugf("fetched 0x%02X\n", op)
	return op
}

// builds instruction out of other instructions
func build(is ...instruction) instruction {
	return func(c *CPU) {
		for _, i := range is {
			i(c)
		}
	}
}

func label(s string) instruction {
	return func(c *CPU) {
		c.Debugf("decode %s\n", s)
	}
}

// decode distinguishes the instruction and reads operands if necessary
// todo: instructions should be created statically at start up
func (c *CPU) decode(b byte) instruction {
	switch b {
	case 0x00:
		return build(
			label("noop"),
			noop,
		)
	case 0x05:
		return build(
			label("DEC B"),
			dec_reg(B),
		)
	case 0x06:
		return build(
			label("LD B, d8"),
			ld_reg_d8(B),
		)
	case 0x0C:
		return build(
			label("INC C"),
			inc_reg(C),
		)
	case 0x0D:
		return build(
			label("DEC C"),
			dec_reg(C),
		)
	case 0x0E:
		return build(
			label("LD C, d8"),
			ld_reg_d8(C),
		)
	case 0x11:
		lsb := c.readByte()
		msb := c.readByte()
		return build(
			label("LD DE, d16"),
			ld_word(D, E, msb, lsb),
		)
	case 0x17:
		return build(
			label("RLA"),
			rl_reg(A),
		)
	case 0x18:
		return build(
			label("JR r8"),
			jr_r8,
		)
	case 0x1A:
		addr := toWord(c.R[D], c.R[E])
		return build(
			label("LD A, (DE)"),
			ld_reg_addr(A, addr),
		)
	case 0x13:
		return build(
			label("INC DE"),
			inc_nn(D, E),
		)
	case 0x20:
		return build(
			label("JR NZ, r8"),
			jr_nz_r8,
		)
	case 0x21:
		lsb := c.readByte()
		msb := c.readByte()
		return build(
			label("ld HL, d16"),
			ld_word(H, L, msb, lsb),
		)
	case 0x22:
		addr := toWord(c.R[H], c.R[L])
		return build(
			label("LD (HL+), A"),
			ldi_addr_reg(addr, A),
		)
	case 0x23:
		return build(
			label("INC HL"),
			inc_nn(H, L),
		)
	case 0x28:
		return build(
			label("JR Z, r8"),
			jr_z_r8,
		)
	case 0x2E:
		return build(
			label("LD L, d8"),
			ld_reg_d8(L),
		)
	case 0x31:
		lsb := c.readByte()
		msb := c.readByte()
		return build(
			label("LD SP, d16"),
			ld_sp_word(msb, lsb),
		)
	case 0x32:
		addr := toWord(c.R[H], c.R[L])
		return build(
			label("LD (HL-), A"),
			ldd_addr_reg(addr, A),
		)
	case 0x3D:
		return build(
			label("DEC A"),
			dec_reg(A),
		)
	case 0x3E:
		return build(
			label("LD A, d8"),
			ld_reg_d8(A),
		)
	case 0x47:
		return build(
			label("LD B, A"),
			ld_reg_reg(B, A),
		)
	case 0x4F:
		return build(
			label("LD C, A"),
			ld_reg_reg(C, A),
		)
	case 0x61:
		return build(
			label("LD H, C"),
			ld_reg_reg(H, C),
		)
	case 0x77:
		addr := toWord(c.R[H], c.R[L])
		return build(
			label("LD (HL), A"),
			ld_a16_reg(addr, A),
		)
	case 0x7B:
		return build(
			label("LD A, E"),
			ld_reg_reg(A, E),
		)
	case 0xA7:
		return build(
			label("AND A"),
			and_reg(A),
		)
	case 0xAF:
		return build(
			label("XOR A"),
			xor(c.R[A]),
		)
	case 0xC1:
		return build(
			label("POP BC"),
			pop(B, C),
		)
	case 0xC5:
		return build(
			label("PUSH BC"),
			push(B, C),
		)
	case 0xC9:
		return build(
			label("RET"),
			ret,
		)
	case 0xCB:
		return c.decodeExtended(c.readByte())
	case 0xCD:
		lsb := c.readByte()
		msb := c.readByte()
		addr := toWord(msb, lsb)
		return build(
			label("CALL a16"),
			call(addr),
		)
	case 0xC3:
		return build(
			label("JMP a16"),
			jmp_a16,
		)
	case 0xE0:
		return build(
			label("LDH (a8), A"),
			ldh_a8_reg(A),
		)
	case 0xE2:
		return build(
			label("LD (C), A"),
			ld_offset_addr(C, A),
		)
	case 0xEA:
		lsb := c.readByte()
		msb := c.readByte()
		addr := toWord(msb, lsb)
		return build(
			label("LD (a16), A"),
			ld_a16_reg(addr, A),
		)
	case 0xF0:
		return build(
			label("LDH A, (a8)"),
			ldh_reg_a8(A),
		)
	case 0xF3:
		return build(
			label("DI"),
			di,
		)
	case 0xFE:
		return build(
			label("CP d8"),
			cp_byte(c.readByte()),
		)
	default:
		return instructionNotImplemented(b)
	}
}

func (c *CPU) decodeExtended(b byte) instruction {
	c.Debugf("fetched extended instruction 0x%02X\n", b)
	switch b {
	case 0x11:
		return build(
			label("RL C"),
			rl_reg(C),
		)
	case 0x7C:
		return bit(7, H)
	default:
		return extendedInstructionNotImplemented(b)
	}
}

// read a byte from the PC a.k.a `n`
func (c *CPU) readByte() uint8 {
	b := c.MMU.ReadByte(c.PC)
	c.PC++
	c.M++
	c.T += 4
	return b
}

// combine two bytes into a word such that a is upper
// half and b is lower half
func toWord(a, b uint8) uint16 {
	return uint16(a)<<8 | uint16(b)
}

type GPU interface {
	GetScrollY() byte
	SetScrollY(y byte)

	GetScrollX() byte
	SetScrollX(x byte)

	SetStat(s byte)
	GetStat() byte

	SetControl(c byte)
	GetControl() byte

	GetLY() byte
	ResetLY()
	// todo: deprecate above methods and just rely on Module interface
	Module

	// Loop sets up graphics and runs the game loop, calling the update function on each tick
	Loop(update func())
}

// Module represents another memory mapped module such as the GPU or APU
type Module interface {
	GetRegister(addr uint16) byte
	SetRegister(addr uint16, b byte)
}

func (c *CPU) stackPush(b byte) {
	c.SP--
	c.MMU.WriteByte(c.SP, b)
}

func (c *CPU) stackPop() byte {
	b := c.MMU.ReadByte(c.SP)
	c.SP++
	return b
}
