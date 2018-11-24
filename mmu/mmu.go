package mmu

type MMU []uint8

func New() MMU {
	return make(MMU, 0x10000)
}

func (m MMU) ReadByte(addr uint16) uint8 {
	return m[addr]
}

func (m MMU) WriteByte(addr uint16, n uint8) {
	m[addr] = n
}
