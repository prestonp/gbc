package gb

import (
	"testing"

	"github.com/prestonp/gbc/pkg/gb/apu"
	"github.com/prestonp/gbc/pkg/gb/gpu"
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
	apu := apu.New()
	gpu := gpu.New()
	mmu := NewMMU(nil, nil, nil, apu)
	cpu := NewCPU(mmu, gpu, false)

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

func TestInc(t *testing.T) {
	cpu := NewCPU(nil, gpu.New(), false)
	cpu.R[C] = 0xF
	inc := inc_reg(C)
	inc(cpu)
	require.NotEqual(t, FlagZero, cpu.R[F]&FlagZero)
	require.Zero(t, cpu.R[F]&FlagSubtract)
	require.Equal(t, FlagHalfCarry, cpu.R[F]&FlagHalfCarry)

	cpu.R[C] = 0xE
	inc(cpu)
	require.NotEqual(t, FlagZero, cpu.R[F]&FlagZero)
	require.Zero(t, cpu.R[F]&FlagSubtract)
	require.Zero(t, cpu.R[F]&FlagHalfCarry)
}

func TestRotate(t *testing.T) {
	t.Run("rl reg", func(t *testing.T) {
		cpu := NewCPU(nil, gpu.New(), false)
		rotate := rl_reg(B)
		{
			cpu.R[B] = 0x80
			cpu.R[F] = FlagCarry

			rotate(cpu)

			require.EqualValues(t, 0x01, cpu.R[B])
			require.EqualValues(t, FlagCarry, cpu.R[F]&FlagCarry, "7th bit should've set carry flag")
			require.EqualValues(t, 1, cpu.R[B]&1, "0th bit should've shifted in previous carry flag")
		}
		{
			cpu.R[B] = 0x7F
			cpu.R[F] = 0

			rotate(cpu)

			require.EqualValues(t, 0xFE, cpu.R[B])
			require.EqualValues(t, 0, cpu.R[F]&FlagCarry, "7th bit clear carry flag")
			require.EqualValues(t, 0, cpu.R[B]&1, "0th bit should've shifted in previous carry flag")
		}
	})
}

func TestSub(t *testing.T) {
	cpu := NewCPU(nil, gpu.New(), false)
	op := sub(B)

	{
		cpu.R[A] = 5
		cpu.R[B] = 5
		op(cpu)
		require.EqualValues(t, FlagZero, cpu.R[F]&FlagZero)
		require.EqualValues(t, FlagSubtract, cpu.R[F]&FlagSubtract)
		require.EqualValues(t, FlagHalfCarry, cpu.R[F]&FlagHalfCarry)
		require.EqualValues(t, FlagCarry, cpu.R[F]&FlagCarry)
		require.EqualValues(t, 0, cpu.R[A])
	}
	{
		cpu.R[A] = 5
		cpu.R[B] = 2
		op(cpu)
		require.EqualValues(t, 0, cpu.R[F]&FlagZero)
		require.EqualValues(t, FlagSubtract, cpu.R[F]&FlagSubtract)
		require.EqualValues(t, FlagHalfCarry, cpu.R[F]&FlagHalfCarry)
		require.EqualValues(t, FlagCarry, cpu.R[F]&FlagCarry)
		require.EqualValues(t, 3, cpu.R[A])
	}
	{
		cpu.R[A] = 5
		cpu.R[B] = 6
		op(cpu)
		require.EqualValues(t, 0, cpu.R[F]&FlagZero)
		require.EqualValues(t, FlagSubtract, cpu.R[F]&FlagSubtract)
		require.EqualValues(t, 0, cpu.R[F]&FlagHalfCarry)
		require.EqualValues(t, 0, cpu.R[F]&FlagCarry)
		require.EqualValues(t, 0xFF, cpu.R[A])
	}
}

func TestAdd(t *testing.T) {
	cpu := NewCPU(nil, gpu.New(), false)
	{
		cpu.R[A] = 0x08
		_add(cpu, 0x08)
		require.EqualValues(t, 0x10, cpu.R[A])
		require.EqualValues(t, 0, cpu.R[F]&FlagZero)
		require.EqualValues(t, 0, cpu.R[F]&FlagSubtract)
		require.EqualValues(t, FlagHalfCarry, cpu.R[F]&FlagHalfCarry)
		require.EqualValues(t, 0, cpu.R[F]&FlagCarry)
	}
	{
		cpu.R[A] = 0xFF
		_add(cpu, 0x01)
		require.EqualValues(t, 0, cpu.R[A])
		require.EqualValues(t, FlagZero, cpu.R[F]&FlagZero)
		require.EqualValues(t, 0, cpu.R[F]&FlagSubtract)
		require.EqualValues(t, FlagHalfCarry, cpu.R[F]&FlagHalfCarry)
		require.EqualValues(t, FlagCarry, cpu.R[F]&FlagCarry)
	}
	{
		cpu.R[A] = 0x80
		_add(cpu, 0x80)
		require.EqualValues(t, 0, cpu.R[A])
		require.EqualValues(t, FlagZero, cpu.R[F]&FlagZero)
		require.EqualValues(t, 0, cpu.R[F]&FlagSubtract)
		require.EqualValues(t, 0, cpu.R[F]&FlagHalfCarry)
		require.EqualValues(t, FlagCarry, cpu.R[F]&FlagCarry)
	}
}

func TestCpl(t *testing.T) {
	cpu := NewCPU(nil, gpu.New(), false)
	{
		cpu.R[A] = 0x0F
		cpl(cpu)
		require.EqualValues(t, 0xF0, cpu.R[A])
		require.EqualValues(t, 0x06, cpu.R[F])
	}
	{
		cpu.R[A] = 0xAA
		cpl(cpu)
		require.EqualValues(t, 0x55, cpu.R[A])
		require.EqualValues(t, 0x06, cpu.R[F])
	}
}

func TestSwap(t *testing.T) {
	cpu := NewCPU(nil, gpu.New(), false)
	swapA := swap_reg(A)
	{
		cpu.R[A] = 0x5F
		swapA(cpu)
		require.EqualValues(t, 0xF5, cpu.R[A])
		require.EqualValues(t, 0, cpu.R[F])
	}
	{
		cpu.R[A] = 0xF5
		swapA(cpu)
		require.EqualValues(t, 0x5F, cpu.R[A])
		require.EqualValues(t, 0, cpu.R[F])
	}
	{
		cpu.R[A] = 0x00
		swapA(cpu)
		require.EqualValues(t, 0x00, cpu.R[A])
		require.EqualValues(t, FlagZero, cpu.R[F])
	}
}
