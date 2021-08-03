package gpu

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBGPalette(t *testing.T) {
	gpu := New()
	gpu.WriteByte(0xFF47, 0xFC)
	require.EqualValues(t, color.White, gpu.getColor(0))
	require.EqualValues(t, color.Black, gpu.getColor(1))
	require.EqualValues(t, color.Black, gpu.getColor(2))
	require.EqualValues(t, color.Black, gpu.getColor(3))

	gpu.WriteByte(0xFF47, 0x1B)
	require.EqualValues(t, color.Black, gpu.getColor(0))
	require.EqualValues(t, color.RGBA{192, 192, 192, 255}, gpu.getColor(1))
	require.EqualValues(t, color.RGBA{128, 128, 128, 255}, gpu.getColor(2))
	require.EqualValues(t, color.White, gpu.getColor(3))
}
