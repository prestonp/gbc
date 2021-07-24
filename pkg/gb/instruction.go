package gb

import "log"

type instruction func(c *CPU)

func instructionNotImplemented(op byte) instruction {
	return func(c *CPU) {
		log.Fatalf("unimplemented instruction 0x%02X", op)
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
		ld_a16_reg(addr, r)(c)
		dec_nn(H, L)(c)
	}
}

func ld_word(H, L Register, msb, lsb byte) instruction {
	return func(c *CPU) {
		c.R[H] = msb
		c.R[L] = lsb
		c.Debugf("exec ld HL 0x%02X%02X\n", msb, lsb)
	}
}

func ld_sp_word(msb, lsb byte) instruction {
	return func(c *CPU) {
		c.SP = toWord(msb, lsb)
		c.Debugf("exec ld SP 0x%02X%02X\n", msb, lsb)
	}
}

func ld_reg_reg(dst, src Register) instruction {
	return func(c *CPU) {
		c.R[dst] = c.R[src]
		c.Debugf("exec ld %s %s\n", dst, src)
	}
}

func ld_reg_d8(dst Register) instruction {
	return func(c *CPU) {
		c.R[dst] = c.readByte()
		c.Debugf("exec ld %s 0x%02X\n", dst, c.R[dst])
	}
}

func ld_a16_reg(addr uint16, reg Register) instruction {
	return func(c *CPU) {
		c.MMU.WriteByte(addr, c.R[reg])
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
func cp_byte(b byte) instruction {
	return func(c *CPU) {
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
}

func dec_a16(addr uint16) instruction {
	return func(c *CPU) {
		panic("reimplement this")
		// c.Debugf("TEMP %s\n", c.String())
		// val := c.MMU.ReadByte(addr) - 1
		// c.MMU.WriteByte(addr, val)

		// c.R[F] = 0
		// if val == 0 {
		// 	c.R[F] |= FlagZero
		// }
		// c.R[F] |= FlagSubtract

		// c.Debugf("todo: implement half carry flag in CP instructions\n")
	}
}

func dec_nn(upper, lower Register) instruction {
	return func(c *CPU) {
		word := toWord(c.R[upper], c.R[lower]) - 1
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
		}
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
	offset := c.readByte()
	c.PC += uint16(offset)
	c.Debugf("exec jr r8 -- jumping to 0x%04X\n", c.PC)
}

func and_reg(r Register) instruction {
	return func(c *CPU) {
		c.R[A] = c.R[r] & c.R[A]
		c.R[F] = 0
		if c.R[A] == 0 {
			c.R[F] |= FlagZero
		}
		c.R[F] |= FlagHalfCarry
		c.Debugf("exec and %s flags = 0b%04b\n", r, c.R[F]>>4)
	}
}

func xor(b byte) instruction {
	return func(c *CPU) {
		c.R[A] ^= b
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
