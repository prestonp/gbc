package cpu

import (
	"log"
	"testing"
)

func TestLD_B(t *testing.T) {
	cpu := New()
	cpu.MMU.WriteByte(0x100, 0x06)
	cpu.MMU.WriteByte(0x101, 0x42)

	cpu.Tick()

	if cpu.R[B] != 0x42 {
		log.Printf("expected %x, got %x\n", 0x42, cpu.R[B])
		t.Fail()
	}

	if cpu.M != 2 {
		log.Printf("expected %x, got %x\n", 2, cpu.M)
		t.Fail()
	}

	if cpu.T != 8 {
		log.Printf("expected %x, got %x\n", 8, cpu.T)
		t.Fail()
	}
}

func TestLD_A_B(t *testing.T) {
	cpu := New()
	cpu.MMU.WriteByte(0x100, 0x78)
	cpu.R[B] = 0x42

	cpu.Tick()

	if cpu.R[A] != 0x42 {
		log.Printf("expected %x, got %x\n", 0x42, cpu.R[A])
		t.Fail()
	}

	if cpu.M != 1 {
		log.Printf("expected %x, got %x\n", 1, cpu.M)
		t.Fail()
	}

	if cpu.T != 4 {
		log.Printf("expected %x, got %x\n", 4, cpu.T)
		t.Fail()
	}
}

func TestLD_A_HL_Addr(t *testing.T) {
	cpu := New()
	cpu.MMU.WriteByte(0x100, 0x7E)
	cpu.R[H] = 0x12
	cpu.R[L] = 0x34
	cpu.MMU.WriteByte(0x1234, 0x42)

	cpu.Tick()

	if cpu.R[A] != 0x42 {
		log.Printf("expected %x, got %x\n", 0x42, cpu.R[A])
		t.Fail()
	}

	if cpu.M != 2 {
		log.Printf("expected %x, got %x\n", 2, cpu.M)
		t.Fail()
	}

	if cpu.T != 8 {
		log.Printf("expected %x, got %x\n", 8, cpu.T)
		t.Fail()
	}
}

func Test_HL_Addr_Byte(t *testing.T) {
	cpu := New()
	cpu.R[H] = 0x12
	cpu.R[L] = 0x34
	cpu.MMU.WriteByte(0x100, 0x36)
	cpu.MMU.WriteByte(0x101, 0x42)

	cpu.Tick()

	b := cpu.MMU.ReadByte(0x1234)
	if b != 0x42 {
		log.Printf("expected %x, got %x\n", 0x42, b)
		t.Fail()
	}

	if cpu.M != 3 {
		log.Printf("expected %x, got %x\n", 4, cpu.M)
		t.Fail()
	}

	if cpu.T != 12 {
		log.Printf("expected %x, got %x\n", 16, cpu.T)
		t.Fail()
	}
}

func TestLD_Word_A(t *testing.T) {
	cpu := New()
	cpu.MMU.WriteByte(0x100, 0xEA)
	cpu.MMU.WriteByte(0x101, 0x12)
	cpu.MMU.WriteByte(0x102, 0x34)
	cpu.R[A] = 0x42

	cpu.Tick()

	b := cpu.MMU.ReadByte(0x1234)
	if b != 0x42 {
		log.Printf("expected %x, got %x\n", 0x42, b)
		t.Fail()
	}

	if cpu.M != 4 {
		log.Printf("expected %x, got %x\n", 4, cpu.M)
		t.Fail()
	}

	if cpu.T != 16 {
		log.Printf("expected %x, got %x\n", 16, cpu.T)
		t.Fail()
	}
}

func TestLD_A_C_Addr(t *testing.T) {
	// LD A, (C)
	cpu := New()
	cpu.R[C] = 0x32
	cpu.MMU.WriteByte(0x100, 0xF2)
	cpu.MMU.WriteByte(0xFF32, 0x42)

	cpu.Tick()

	b := cpu.MMU.ReadByte(0xFF32)
	if b != 0x42 {
		log.Printf("expected %x, got %x\n", 0x42, b)
		t.Fail()
	}

	if cpu.M != 2 {
		log.Printf("expected %x, got %x\n", 2, cpu.M)
		t.Fail()
	}

	if cpu.T != 8 {
		log.Printf("expected %x, got %x\n", 8, cpu.T)
		t.Fail()
	}
}

func TestLDD_A_HL(t *testing.T) {
	// LDD A, (HL)
	cpu := New()
	cpu.R[H] = 0x12
	cpu.R[L] = 0x34
	cpu.MMU.WriteByte(0x100, 0x3A)
	cpu.MMU.WriteByte(0x1234, 0x42)

	cpu.Tick()

	a := cpu.R[A]
	if a != 0x42 {
		log.Printf("expected %x, got %x\n", 0x42, a)
		t.Fail()
	}

	// should decrement (HL)
	hl := toWord(cpu.R[H], cpu.R[L])
	if hl != 0x1233 {
		log.Printf("expected %x, got %x\n", 0x1233, hl)
		t.Fail()
	}

	if cpu.M != 2 {
		log.Printf("expected %x, got %x\n", 2, cpu.M)
		t.Fail()
	}

	if cpu.T != 8 {
		log.Printf("expected %x, got %x\n", 8, cpu.T)
		t.Fail()
	}
}

func TestLDH_byte_A(t *testing.T) {
	// LDH (n), A
	cpu := New()
	cpu.R[A] = 0x42
	cpu.MMU.WriteByte(0x100, 0xE0)
	cpu.MMU.WriteByte(0x101, 0x01)

	cpu.Tick()

	b := cpu.MMU.ReadByte(0xFF01)
	if b != 0x42 {
		log.Printf("expected %x, got %x\n", 0x42, b)
		t.Fail()
	}

	if cpu.M != 3 {
		log.Printf("expected %x, got %x\n", 3, cpu.M)
		t.Fail()
	}

	if cpu.T != 12 {
		log.Printf("expected %x, got %x\n", 12, cpu.T)
		t.Fail()
	}
}
func TestLDH_A_byteAddr(t *testing.T) {
	// LDH (n), A
	cpu := New()
	cpu.MMU.WriteByte(0x100, 0xF0)
	cpu.MMU.WriteByte(0x101, 0x01)
	cpu.MMU.WriteByte(0xFF01, 0x42)

	cpu.Tick()

	b := cpu.R[A]
	if b != 0x42 {
		log.Printf("expected %x, got %x\n", 0x42, b)
		t.Fail()
	}

	if cpu.M != 3 {
		log.Printf("expected %x, got %x\n", 3, cpu.M)
		t.Fail()
	}

	if cpu.T != 12 {
		log.Printf("expected %x, got %x\n", 12, cpu.T)
		t.Fail()
	}
}

func TestToWord(t *testing.T) {
	w := toWord(0x12, 0x34)
	if w != 0x1234 {
		log.Printf("expected %x, got %x\n", 0x1214, w)
		t.Fail()
	}
}
