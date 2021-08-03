package gb

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	GPU Module

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

func NewCPU(mmu *MMU, gpu Module, debug bool) *CPU {
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
	done := make(chan bool)
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				c.Update()
			}
		}
	}()

	c.GPU.Run()
}

func (c *CPU) Update() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(c.log.String())
			fmt.Println(r)
			os.Exit(1)
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

	// todo: remove this once timing is implemented
	if op == 0x04 {
		c.Debugf("debugging breakpoint to let gpu render a frame")
		time.Sleep(1 * time.Minute)
	}
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

var extendedOps = map[byte]instruction{
	0x11: build(label("RL C"), rl_reg(C)),
	0x7C: bit(7, H),
}

var ops = map[byte]instruction{
	0x00: build(label("noop"), noop),
	0x05: build(label("DEC B"), dec_reg(B)),
	0x06: build(label("LD B, d8"), ld_reg_d8(B)),
	0x0C: build(label("INC C"), inc_reg(C)),
	0x0D: build(label("DEC C"), dec_reg(C)),
	0x0E: build(label("LD C, d8"), ld_reg_d8(C)),

	0x11: build(label("LD DE, d16"), ld_word(D, E)),
	0x17: build(label("RLA"), rl_reg(A)),
	0x18: build(label("JR r8"), jr_r8),
	0x1A: build(label("LD A, (DE)"), ld_reg_word(A, D, E)),
	0x13: build(label("INC DE"), inc_nn(D, E)),

	0x20: build(label("JR NZ, r8"), jr_nz_r8),
	0x21: build(label("ld HL, d16"), ld_word(H, L)),
	0x22: build(label("LD (HL+), A"), ldi_hl_reg(A)),
	0x23: build(label("INC HL"), inc_nn(H, L)),
	0x28: build(label("JR Z, r8"), jr_z_r8),
	0x2E: build(label("LD L, d8"), ld_reg_d8(L)),

	0x31: build(label("LD SP, d16"), ld_sp_word),
	0x32: build(label("LD (HL-), A"), ldd_hl_reg(A)),
	0x3D: build(label("DEC A"), dec_reg(A)),
	0x3E: build(label("LD A, d8"), ld_reg_d8(A)),

	0x47: build(label("LD B, A"), ld_reg_reg(B, A)),
	0x4F: build(label("LD C, A"), ld_reg_reg(C, A)),

	0x57: build(label("LD D, A"), ld_reg_reg(D, A)),

	0x61: build(label("LD H, C"), ld_reg_reg(H, C)),
	0x67: build(label("LD H, A"), ld_reg_reg(H, A)),

	0x77: build(label("LD (HL), A"), ld_addrhl_reg(A)),
	0x7B: build(label("LD A, E"), ld_reg_reg(A, E)),

	0xA7: build(label("AND A"), and_reg(A)),
	0xAF: build(label("XOR A"), xor_reg(A)),

	0xC1: build(label("POP BC"), pop(B, C)),
	0xC5: build(label("PUSH BC"), push(B, C)),
	0xC9: build(label("RET"), ret),
	0xCD: build(label("CALL a16"), call_a16),
	0xC3: build(label("JMP a16"), jmp_a16),

	0xE0: build(label("LDH (a8), A"), ldh_a8_reg(A)),
	0xE2: build(label("LD (C), A"), ld_offset_addr(C, A)),
	0xEA: build(label("LD (a16), A"), ld_a16_reg(A)),

	0xF0: build(label("LDH A, (a8)"), ldh_reg_a8(A)),
	0xF3: build(label("DI"), di),
	0xFE: build(label("CP d8"), cp_byte),
}

// decode distinguishes the instructions
func (c *CPU) decode(b byte) instruction {
	if b == 0xCB {
		return c.decodeExtended(c.readByte())
	}

	if op, ok := ops[b]; !ok {
		return instructionNotImplemented(b)
	} else {
		return op
	}
}

func (c *CPU) decodeExtended(b byte) instruction {
	c.Debugf("fetched extended instruction 0x%02X\n", b)
	if op, ok := extendedOps[b]; !ok {
		return extendedInstructionNotImplemented(b)
	} else {
		return op
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

// Module represents another memory mapped module such as the GPU or APU
type Module interface {
	ReadByte(addr uint16) byte
	WriteByte(addr uint16, b byte)
	Run()
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
