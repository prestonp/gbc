package logbuf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogBuf(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		buf := New(10)
		msg := []byte("hello")
		n, err := buf.Write(msg)
		require.NoError(t, err)
		require.Equal(t, len(msg), n)
		require.Equal(t, msg, buf.Bytes())
	})
	t.Run("wrap", func(t *testing.T) {
		buf := New(10)
		msg := []byte("hello world")
		n, err := buf.Write(msg)
		require.NoError(t, err)
		require.Equal(t, len(msg), n)
		require.Equal(t, []byte("ello world"), buf.Bytes())
	})
	t.Run("string", func(t *testing.T) {
		buf := New(10)
		msg := []byte("hello world")
		buf.Write(msg)
		require.Equal(t, "ello world", buf.String())
	})
}
