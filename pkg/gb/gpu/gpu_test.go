package gpu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDot(t *testing.T) {
	gpu := New()
	gpu.WriteByte(0xFF47, 0xFC)

	require.EqualValues(t, 0b11, gpu.getPalette(DotDarker))
	require.EqualValues(t, 0b11, gpu.getPalette(DotDark))
	require.EqualValues(t, 0b11, gpu.getPalette(DotLight))
	require.EqualValues(t, 0b00, gpu.getPalette(DotLighter))
}
