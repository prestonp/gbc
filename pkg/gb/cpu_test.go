package gb

import (
	"github.com/prestonp/gbc/pkg/gb/apu"
	"github.com/prestonp/gbc/pkg/gb/gpu"
)

func getTestCPU() *CPU {
	gpu := gpu.New()
	apu := apu.New()
	mmu := NewMMU(nil, nil, gpu, apu)
	return NewCPU(mmu, false)
}

// func TestLD_B(t *testing.T) {
// 	cpu := getTestCPU()
// 	cpu.MMU.WriteByte(0x100, 0x06)
// 	cpu.MMU.WriteByte(0x101, 0x42)

// 	cpu.Tick()

// 	if cpu.R[B] != 0x42 {
// 		log.Printf("expected %x, got %x\n", 0x42, cpu.R[B])
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 2, 8)
// }

// func TestLD_A_B(t *testing.T) {
// 	cpu := getTestCPU()
// 	cpu.MMU.WriteByte(0x100, 0x78)
// 	cpu.R[B] = 0x42

// 	cpu.Tick()

// 	if cpu.R[A] != 0x42 {
// 		log.Printf("expected %x, got %x\n", 0x42, cpu.R[A])
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 1, 4)
// }

// func TestLD_A_HL_Addr(t *testing.T) {
// 	cpu := getTestCPU()
// 	cpu.MMU.WriteByte(0x100, 0x7E)
// 	cpu.R[H] = 0x12
// 	cpu.R[L] = 0x34
// 	cpu.MMU.WriteByte(0x1234, 0x42)

// 	cpu.Tick()

// 	if cpu.R[A] != 0x42 {
// 		log.Printf("expected %x, got %x\n", 0x42, cpu.R[A])
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 2, 8)
// }

// func Test_HL_Addr_Byte(t *testing.T) {
// 	cpu := getTestCPU()
// 	cpu.R[H] = 0x12
// 	cpu.R[L] = 0x34
// 	cpu.MMU.WriteByte(0x100, 0x36)
// 	cpu.MMU.WriteByte(0x101, 0x42)

// 	cpu.Tick()

// 	b := cpu.MMU.ReadByte(0x1234)
// 	if b != 0x42 {
// 		log.Printf("expected %x, got %x\n", 0x42, b)
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 3, 12)
// }

// func TestLD_Word_A(t *testing.T) {
// 	cpu := getTestCPU()
// 	cpu.MMU.WriteByte(0x100, 0xEA)
// 	cpu.MMU.WriteByte(0x101, 0x12)
// 	cpu.MMU.WriteByte(0x102, 0x34)
// 	cpu.R[A] = 0x42

// 	cpu.Tick()

// 	b := cpu.MMU.ReadByte(0x1234)
// 	if b != 0x42 {
// 		log.Printf("expected %x, got %x\n", 0x42, b)
// 		t.Fail()
// 	}
// 	assertTiming(t, cpu, 4, 16)
// }

// func TestLD_A_C_Addr(t *testing.T) {
// 	// LD A, (C)
// 	cpu := getTestCPU()
// 	cpu.R[C] = 0x32
// 	cpu.MMU.WriteByte(0x100, 0xF2)
// 	cpu.MMU.WriteByte(0xFF32, 0x42)

// 	cpu.Tick()

// 	b := cpu.MMU.ReadByte(0xFF32)
// 	if b != 0x42 {
// 		log.Printf("expected %x, got %x\n", 0x42, b)
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 2, 8)
// }

// func TestLDD_A_HL(t *testing.T) {
// 	// LDD A, (HL)
// 	cpu := getTestCPU()
// 	cpu.R[H] = 0x12
// 	cpu.R[L] = 0x34
// 	cpu.MMU.WriteByte(0x100, 0x3A)
// 	cpu.MMU.WriteByte(0x1234, 0x42)

// 	cpu.Tick()

// 	a := cpu.R[A]
// 	if a != 0x42 {
// 		log.Printf("expected %x, got %x\n", 0x42, a)
// 		t.Fail()
// 	}

// 	// should decrement (HL)
// 	hl := toWord(cpu.R[H], cpu.R[L])
// 	if hl != 0x1233 {
// 		log.Printf("expected %x, got %x\n", 0x1233, hl)
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 2, 8)
// }

// func TestLDH_byte_A(t *testing.T) {
// 	// LDH (n), A
// 	cpu := getTestCPU()
// 	cpu.R[A] = 0x42
// 	cpu.MMU.WriteByte(0x100, 0xE0)
// 	cpu.MMU.WriteByte(0x101, 0x01)

// 	cpu.Tick()

// 	b := cpu.MMU.ReadByte(0xFF01)
// 	if b != 0x42 {
// 		log.Printf("expected %x, got %x\n", 0x42, b)
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 3, 12)
// }
// func TestLDH_A_byteAddr(t *testing.T) {
// 	// LDH (n), A
// 	cpu := getTestCPU()
// 	cpu.MMU.WriteByte(0x100, 0xF0)
// 	cpu.MMU.WriteByte(0x101, 0x01)
// 	cpu.MMU.WriteByte(0xFF01, 0x42)

// 	cpu.Tick()

// 	b := cpu.R[A]
// 	if b != 0x42 {
// 		log.Printf("expected %x, got %x\n", 0x42, b)
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 3, 12)
// }

// func TestLD_word_word(t *testing.T) {
// 	cpu := getTestCPU()
// 	cpu.MMU.WriteByte(0x100, 0x01)
// 	cpu.MMU.WriteByte(0x101, 0x12)
// 	cpu.MMU.WriteByte(0x102, 0x34)

// 	cpu.Tick()

// 	if cpu.R[B] != 0x12 {
// 		log.Printf("expected %x, got %x\n", 0x12, cpu.R[B])
// 		t.Fail()
// 	}

// 	if cpu.R[C] != 0x34 {
// 		log.Printf("expected %x, got %x\n", 0x34, cpu.R[C])
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 3, 12)
// }

// func TestLD_SP_nn(t *testing.T) {
// 	cpu := getTestCPU()
// 	cpu.MMU.WriteByte(0x100, 0x31)
// 	cpu.MMU.WriteByte(0x101, 0x12)
// 	cpu.MMU.WriteByte(0x102, 0x34)

// 	cpu.Tick()

// 	if cpu.SP != 0x1234 {
// 		log.Printf("expected %x, got %x\n", 0x1234, cpu.SP)
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 3, 12)
// }

// func TestLD_SP_HL(t *testing.T) {
// 	cpu := getTestCPU()
// 	cpu.MMU.WriteByte(0x100, 0xF9)
// 	cpu.R[H] = 0x12
// 	cpu.R[L] = 0x34

// 	cpu.Tick()

// 	if cpu.SP != 0x1234 {
// 		log.Printf("expected %x, got %x\n", 0x1234, cpu.SP)
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 2, 8)
// }

// func TestLD_HL_SP_n(t *testing.T) {
// 	cpu := getTestCPU()
// 	cpu.MMU.WriteByte(0x100, 0xF8)
// 	cpu.MMU.WriteByte(0x101, 0x01)
// 	cpu.SP = 0x1234

// 	cpu.Tick()

// 	if cpu.R[H] != 0x12 {
// 		log.Printf("expected %x, got %x\n", 0x12, cpu.R[H])
// 		t.Fail()
// 	}

// 	if cpu.R[L] != 0x35 {
// 		log.Printf("expected %x, got %x\n", 0x35, cpu.R[L])
// 		t.Fail()
// 	}

// 	if cpu.R[F] != 0x0 {
// 		log.Printf("expected %x, got %x\n", 0x0, cpu.R[F])
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 3, 12)
// }

// func TestLD_nn_SP(t *testing.T) {
// 	cpu := getTestCPU()
// 	cpu.MMU.WriteByte(0x100, 0x08)
// 	cpu.MMU.WriteByte(0x101, 0x12)
// 	cpu.MMU.WriteByte(0x102, 0x34)
// 	cpu.SP = 0x1738

// 	cpu.Tick()

// 	if cpu.MMU.ReadByte(0x1234) != 0x38 {
// 		log.Printf("expected %x, got %x\n", 0x38, cpu.MMU.ReadByte(0x1234))
// 		t.Fail()
// 	}

// 	if cpu.MMU.ReadByte(0x1235) != 0x17 {
// 		log.Printf("expected %x, got %x\n", 0x17, cpu.MMU.ReadByte(0x1235))
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 5, 20)
// }

// func TestStack(t *testing.T) {
// 	cpu := getTestCPU()

// 	// Test push
// 	cpu.MMU.WriteByte(0x100, 0xF6)
// 	cpu.R[A] = 0x12
// 	cpu.R[F] = 0x34
// 	cpu.SP = 0x5555

// 	cpu.Tick()

// 	if cpu.MMU.ReadByte(0x5555) != cpu.R[F] {
// 		log.Printf("expected %x, got %x\n", cpu.R[F], cpu.MMU.ReadByte(0x5555))
// 		t.Fail()
// 	}

// 	if cpu.MMU.ReadByte(0x5554) != cpu.R[A] {
// 		log.Printf("expected %x, got %x\n", cpu.R[A], cpu.MMU.ReadByte(0x5554))
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 4, 16)

// 	// Test pop
// 	cpu.MMU.WriteByte(0x101, 0xF1)

// 	// Replace AF with other values
// 	cpu.R[A] = 0xFF
// 	cpu.R[F] = 0xFF

// 	// Reset timing for checking pop timing
// 	cpu.M = 0
// 	cpu.T = 0

// 	cpu.Tick()

// 	if cpu.R[A] != 0x12 {
// 		log.Printf("expected %x, got %x\n", 0x12, cpu.R[A])
// 		t.Fail()
// 	}

// 	if cpu.R[F] != 0x34 {
// 		log.Printf("expected %x, got %x\n", 0x34, cpu.R[A])
// 		t.Fail()
// 	}

// 	assertTiming(t, cpu, 3, 12)
// }

// func assertTiming(t *testing.T, cpu *CPU, mCycles, tCycles int) {
// 	if cpu.M != mCycles {
// 		log.Printf("expected %d, got %d\n", mCycles, cpu.M)
// 		t.Fail()
// 	}

// 	if cpu.T != tCycles {
// 		log.Printf("expected %d, got %d\n", tCycles, cpu.T)
// 		t.Fail()
// 	}
// }

// func TestToWord(t *testing.T) {
// 	w := toWord(0x12, 0x34)
// 	if w != 0x1234 {
// 		log.Printf("expected %x, got %x\n", 0x1214, w)
// 		t.Fail()
// 	}
// }
