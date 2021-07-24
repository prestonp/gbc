package gb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMMU(t *testing.T) {
	rom, err := ReadRom("testdata/tetris.gb")
	require.NoError(t, err)
	require.NotEmpty(t, rom)
}
