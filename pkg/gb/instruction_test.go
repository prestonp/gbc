package gb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddSignedByte(t *testing.T) {
	t.Run("add positive int8", func(t *testing.T) {
		require.Equal(t, uint16(108), addSignedByte(uint16(101), int8(7)))
	})
	t.Run("add negative int8", func(t *testing.T) {
		require.Equal(t, uint16(94), addSignedByte(uint16(101), int8(-7)))
	})
	t.Run("add 0", func(t *testing.T) {
		require.Equal(t, uint16(101), addSignedByte(uint16(101), 0))
	})
}

func TestBit(t *testing.T) {
	mmu := NewMMU(nil, nil, nil)
	cpu := NewCPU(mmu, false)

	t.Run("check specific bit in a register", func(t *testing.T) {
		cpu.R[H] = 0x80
		check := bit(7, H)
		check(cpu)
		require.Zero(t, cpu.R[F]&FlagZero)
		require.Zero(t, cpu.R[F]&FlagSubtract)
		require.Equal(t, cpu.R[F]&FlagHalfCarry, FlagHalfCarry)

		cpu.R[H] = 0x00
		check(cpu)
		require.Equal(t, cpu.R[F]&FlagZero, FlagZero)
		require.Zero(t, cpu.R[F]&FlagSubtract)
		require.Equal(t, cpu.R[F]&FlagHalfCarry, FlagHalfCarry)
	})
}
