package apu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSweep(t *testing.T) {
	apu := New()
	apu.SetRegister(0xFF10, 0xFF)
	require.EqualValues(t, 0b111, apu.sweepTime())
	require.EqualValues(t, true, apu.sweepMode())
	require.EqualValues(t, 0b111, apu.sweepShift())

	apu.SetRegister(0xFF10, 0xF3)
	require.EqualValues(t, 0b111, apu.sweepTime())
	require.EqualValues(t, false, apu.sweepMode())
	require.EqualValues(t, 0b011, apu.sweepShift())

	apu.SetRegister(0xFF10, 0x00)
	require.EqualValues(t, 0b000, apu.sweepTime())
	require.EqualValues(t, false, apu.sweepMode())
	require.EqualValues(t, 0, apu.sweepShift())
}
