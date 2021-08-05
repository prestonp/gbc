package gb

import "log"

type instruction func(c *CPU)

func instructionNotImplemented(op byte) instruction {
	return func(c *CPU) {
		log.Panicf("unimplemented instruction 0x%02X", op)
	}
}
func extendedInstructionNotImplemented(op byte) instruction {
	return func(c *CPU) {
		log.Panicf("unimplemented extended instruction 0xCB 0x%02X", op)
	}
}

func noop(c *CPU) {
	c.Debugf("noop\n")
}

func jmp_a16(c *CPU) {
	lsb := c.readByte()
	msb := c.readByte()
	addr := toWord(msb, lsb)
	c.PC = addr
	c.Debugf("exec jumped to 0x%04X\n", addr)
}

func ldd_addr_reg(addr uint16, r Register) instruction {
	return func(c *CPU) {
		ld_addr_reg(addr, r)(c)
		dec_nn(H, L)(c)
	}
}

func ldd_hl_reg(r Register) instruction {
	return func(c *CPU) {
		addr := toWord(c.R[H], c.R[L])
		ldd := ldd_addr_reg(addr, r)
		ldd(c)
	}
}

func ldi_addr_reg(addr uint16, r Register) instruction {
	return func(c *CPU) {
		ld_addr_reg(addr, r)(c)
		inc_nn(H, L)(c)
	}
}

func ldi_hl_reg(r Register) instruction {
	return func(c *CPU) {
		addr := toWord(c.R[H], c.R[L])
		ldi := ldi_addr_reg(addr, r)
		ldi(c)
	}
}

func ldi_reg_word(dst, upper, lower Register) instruction {
	return func(c *CPU) {
		ld_reg_word(dst, upper, lower)(c)
		inc_nn(upper, lower)(c)
	}
}

func ld_word(upper, lower Register) instruction {
	return func(c *CPU) {
		c.R[lower] = c.readByte()
		c.R[upper] = c.readByte()
		c.Debugf("exec ld %s%s 0x%02X%02X\n", upper, lower, c.R[upper], c.R[lower])
	}
}

func ld_addrhl_reg(r Register) instruction {
	return func(c *CPU) {
		addr := toWord(c.R[H], c.R[L])
		c.MMU.WriteByte(addr, c.R[r])
	}
}

func ld_addrhl_d8(c *CPU) {
	addr := toWord(c.R[H], c.R[L])
	c.MMU.WriteByte(addr, c.readByte())
}

func ld_sp_word(c *CPU) {
	lsb := c.readByte()
	msb := c.readByte()
	c.SP = toWord(msb, lsb)
	c.Debugf("exec ld SP 0x%02X%02X\n", msb, lsb)
}

func ld_reg_reg(dst, src Register) instruction {
	return func(c *CPU) {
		c.R[dst] = c.R[src]
		c.Debugf("exec LD %s 0x%02X\n", dst, c.R[src])
	}
}

func ld_reg_d8(dst Register) instruction {
	return func(c *CPU) {
		c.R[dst] = c.readByte()
		c.Debugf("exec ld %s 0x%02X\n", dst, c.R[dst])
	}
}

func ld_reg_addr(r Register, addr uint16) instruction {
	return func(c *CPU) {
		c.R[r] = c.MMU.ReadByte(addr)
		c.Debugf("exec LD %s, 0x%02X\n", r, c.R[r])
	}
}

func ld_reg_word(r, upper, lower Register) instruction {
	return func(c *CPU) {
		addr := toWord(c.R[upper], c.R[lower])
		c.R[r] = c.MMU.ReadByte(addr)
	}
}

func ld_a16_reg(reg Register) instruction {
	return func(c *CPU) {
		lsb := c.readByte()
		msb := c.readByte()
		addr := toWord(msb, lsb)
		ld_addr_reg(addr, reg)(c)
	}
}

func ld_addr_reg(addr uint16, reg Register) instruction {
	return func(c *CPU) {
		c.MMU.WriteByte(addr, c.R[reg])
		c.Debugf("exec ld (0x%04X) = 0x%02X\n", addr, c.R[reg])
	}
}

func ld_offset_addr(offset, src Register) instruction {
	return func(c *CPU) {
		addr := 0xFF00 + uint16(c.R[offset])
		c.MMU.WriteByte(addr, c.R[src])
		c.Debugf("exec ld (0x%04X) = 0x%02X\n", addr, c.R[src])
	}
}

func ldh_a8_reg(src Register) instruction {
	return func(c *CPU) {
		offset := c.readByte()
		addr := 0xFF00 + uint16(offset)
		c.MMU.WriteByte(addr, c.R[src])
		c.Debugf("exec ldh (0x%04X) %s = 0x%02X\n", addr, src, c.R[src])
	}
}

func ldh_reg_a8(dest Register) instruction {
	return func(c *CPU) {
		offset := c.readByte()
		addr := 0xFF00 + uint16(offset)
		c.R[dest] = c.MMU.ReadByte(addr)
		c.Debugf("exec ldh %s %02X\n", dest, c.R[dest])
	}
}

// helper func to compare register A to a byte
func _compare(b byte, c *CPU) {
	c.R[F] = 0
	if c.R[A] == b {
		c.R[F] |= FlagZero
	} else if c.R[A] < b {
		c.R[F] |= FlagCarry
	}
	var diff byte = c.R[A] - b
	if diff&0xF > c.R[A]&0xF {
		c.R[F] |= FlagHalfCarry
	}
	c.R[F] |= FlagSubtract
	c.Debugf("exec cp d8 (0x%02X) flags = 0b%04b\n", b, c.R[F]>>4)
}

func cp_byte(c *CPU) {
	_compare(c.readByte(), c)
}

func cp_hl(c *CPU) {
	addr := toWord(c.R[H], c.R[L])
	_compare(c.MMU.ReadByte(addr), c)
}

func dec_nn(upper, lower Register) instruction {
	return func(c *CPU) {
		word := toWord(c.R[upper], c.R[lower]) - 1
		c.R[upper] = byte(word >> 8)
		c.R[lower] = byte(word & 0xFF)
	}
}
func inc_nn(upper, lower Register) instruction {
	return func(c *CPU) {
		word := toWord(c.R[upper], c.R[lower]) + 1
		c.R[upper] = byte(word >> 8)
		c.R[lower] = byte(word & 0xFF)
	}
}

func dec_reg(r Register) instruction {
	return func(c *CPU) {
		c.R[r]--

		c.R[F] = 0
		if c.R[r] == 0 {
			c.R[F] |= FlagZero
		}
		c.R[F] |= FlagSubtract
		if c.R[F] == 0xF {
			c.R[F] |= FlagHalfCarry
		} else {
			c.R[F] &= ^FlagHalfCarry
		}
	}
}

func halfCarryAdd(a, b byte) bool {
	return (a&0xF+b&0xF)&0x10 == 0x10
}

func fullCarryAdd(a, b byte) bool {
	return (uint16(a)+uint16(b))&0x100 == 0x100
}

func inc_reg(r Register) instruction {
	return func(c *CPU) {
		if halfCarryAdd(c.R[r], 1) {
			c.R[F] |= FlagHalfCarry
		} else {
			c.R[F] &= ^FlagHalfCarry
		}

		c.R[r]++

		if c.R[r] == 0 {
			c.R[F] |= FlagZero
		}

		c.R[F] &= ^FlagSubtract
	}
}

func jr_z_r8(c *CPU) {
	offset := int8(c.readByte())
	if c.R[F]&FlagZero > 0 {
		c.PC = addSignedByte(c.PC, offset)
		c.Debugf("exec jr Z r8 -- Z is set, jumping to 0x%04X\n", c.PC)
	} else {
		c.Debugf("exec jr Z r8 -- Z is clear, skipping jump\n")
	}
}

func addSignedByte(val uint16, offset int8) uint16 {
	if offset < 0 {
		return val - uint16(offset*-1)
	}
	return val + uint16(offset)
}

func jr_nz_r8(c *CPU) {
	offset := int8(c.readByte())
	if c.R[F]&FlagZero == 0 {
		c.PC = addSignedByte(c.PC, offset)
		c.Debugf("exec jr NZ r8 -- Z is clear, jumping to 0x%04X\n", c.PC)
	} else {
		c.Debugf("exec jr NZ r8 -- Z is set, skipping jump\n")
	}
}

func jr_r8(c *CPU) {
	offset := int8(c.readByte())
	if offset >= 0 {
		c.PC += uint16(offset)
	} else {
		c.PC -= uint16(offset * -1)
	}
	c.Debugf("exec jr r8 -- jumping to 0x%04X\n", c.PC)
}

func _and(c *CPU, b byte) {
	c.R[A] = b & c.R[A]
	c.R[F] = 0
	if c.R[A] == 0 {
		c.R[F] |= FlagZero
	}
	c.R[F] |= FlagHalfCarry
	c.Debugf("exec AND 0x%02X flags = 0b%04b\n", b, c.R[F]>>4)
}

func and_reg(r Register) instruction {
	return func(c *CPU) {
		_and(c, c.R[r])
	}
}

func and_d8(c *CPU) {
	_and(c, c.readByte())
}

func xor_reg(r Register) instruction {
	return func(c *CPU) {
		c.R[A] ^= c.R[r]
		c.R[F] = 0
		if c.R[A] == 0 {
			c.R[F] |= FlagZero
		}
	}
}

// disable interrupts
func di(c *CPU) {
	c.shouldDI = true
}

// enable interrupts
func ei(c *CPU) {
	c.shouldEI = true
}

func bit(idx int, r Register) instruction {
	return func(c *CPU) {
		c.R[F] = 0
		if (c.R[r] & (1 << idx)) == 0 {
			c.R[F] |= FlagZero
		}
		c.R[F] |= FlagHalfCarry
		c.Debugf("checking bit %d of register %s\n", idx, r)
	}
}

func call_a16(c *CPU) {
	addr := func() uint16 {
		lsb := c.readByte()
		msb := c.readByte()
		return toWord(msb, lsb)
	}()
	c.Debugf("exec call: push PC 0x%04X onto stack, jumping to 0x%04X\n", c.PC, addr)
	lsb := byte(c.PC & 0xFF)
	msb := byte(c.PC >> 8)
	c.stackPush(lsb)
	c.stackPush(msb)
	c.PC = addr
}

func push(upper, lower Register) instruction {
	return func(c *CPU) {
		c.Debugf("exec push 0x%04X onto stack\n", toWord(c.R[upper], c.R[lower]))
		c.stackPush(c.R[lower])
		c.stackPush(c.R[upper])
	}
}

func pop(upper, lower Register) instruction {
	return func(c *CPU) {
		msb := c.stackPop()
		lsb := c.stackPop()
		c.R[upper] = msb
		c.R[lower] = lsb
		c.Debugf("exec POP %s%s = 0x%04X\n", upper, lower, toWord(msb, lsb))
	}
}

func rl_reg(r Register) instruction {
	return func(c *CPU) {
		prevCarry := c.R[F] & FlagCarry
		b7 := c.R[r] & 0x80

		c.R[r] = c.R[r] << 1

		if prevCarry == FlagCarry {
			c.R[r] |= 1
		}

		c.R[F] = 0
		if c.R[r] == 0 {
			c.R[F] |= FlagZero
		}

		if b7 == 0x80 {
			c.R[F] |= FlagCarry
		}
	}
}

func rl_addr(addr uint16) instruction {
	return func(c *CPU) {

	}
}

func ret(c *CPU) {
	msb := c.stackPop()
	lsb := c.stackPop()
	c.PC = toWord(msb, lsb)
	c.Debugf("exec RET - PC jumping to 0x%04X\n", c.PC)
}

func sub(r Register) instruction {
	return func(c *CPU) {
		c.R[F] = 0

		diff := c.R[A] - c.R[r]
		if diff == 0 {
			c.R[F] |= FlagZero
		}

		c.R[F] |= FlagSubtract

		if diff&0xF <= c.R[A]&0xF {
			c.R[F] |= FlagHalfCarry
		}

		if int(c.R[A])-int(c.R[r]) >= 0 {
			c.R[F] |= FlagCarry
		}

		c.R[A] = diff
	}
}

// helper func to add a byte to register A
func _add(c *CPU, b byte) {
	sum := c.R[A] + b

	c.R[F] = 0
	if sum == 0 {
		c.R[F] |= FlagZero
	}

	if halfCarryAdd(c.R[A], b) {
		c.R[F] |= FlagHalfCarry
	}

	if fullCarryAdd(c.R[A], b) {
		c.R[F] |= FlagCarry
	}

	c.R[A] = sum
}

func add_reg(r Register) instruction {
	return func(c *CPU) {
		_add(c, c.R[r])
	}
}

func add_hl(c *CPU) {
	addr := toWord(c.R[H], c.R[L])
	_add(c, c.MMU.ReadByte(addr))
}

func add_d8(c *CPU) {
	_add(c, c.readByte())
}

func or_reg(r Register) instruction {
	return func(c *CPU) {
		c.R[A] |= c.R[r]
		if c.R[A] == 0 {
			c.R[F] = FlagZero
		} else {
			c.R[F] = 0
		}
	}
}

func cpl(c *CPU) {
	c.R[A] = ^c.R[A]
	c.R[F] = 0x06
}

func swap_reg(r Register) instruction {
	return func(c *CPU) {
		c.R[r] = (c.R[r] & 0x0F << 4) | (c.R[r] & 0xF0 >> 4)
		if c.R[r] == 0 {
			c.R[F] = FlagZero
		} else {
			c.R[F] = 0
		}
	}
}

func rst(offset byte) instruction {
	return func(c *CPU) {
		addr := uint16(offset)
		c.Debugf("exec rst 0x%02X: push PC 0x%04X onto stack, jumping to 0x%04X\n", offset, c.PC, addr)
		lsb := byte(c.PC & 0xFF)
		msb := byte(c.PC >> 8)
		c.stackPush(lsb)
		c.stackPush(msb)
		c.PC = addr
	}
}
