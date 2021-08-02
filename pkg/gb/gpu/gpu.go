package gpu

import (
	"fmt"
	"log"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type GPU struct {
	vram []byte
	scx  byte
	scy  byte
	stat byte

	ly byte // lcdc y-coordinate

	// lcd control
	lcdEnable              bool
	winTileMapArea         bool
	winEnable              bool
	bgAndWinTileDataArea   bool
	bgTileMapArea          bool
	objSize                bool
	objEnable              bool
	bgAndWinEnablePriority bool

	bgp byte
}

func (g *GPU) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "\tscx: %d\n", g.scx)
	fmt.Fprintf(&b, "\tscy: %d\n", g.scy)
	fmt.Fprintf(&b, "\tlcdc y: %d\n", g.ly)
	fmt.Fprintf(&b, "\tlcd status: %08b\n", g.stat)
	fmt.Fprintf(&b, "\tlcd control: %08b\n", g.getControl())
	return b.String()
}

func (g *GPU) WriteByte(addr uint16, b byte) {
	switch {
	case addr >= 0x8000 && addr <= 0x9FFF:
		g.vram[addr-0x8000] = b
	case addr == 0xFF40:
		// LCD Control
		g.setControl(b)
	case addr == 0xFF41:
		g.setStat(b & 0xF8)
	case addr == 0xFF42:
		g.setScrollY(b)
	case addr == 0xFF43:
		g.setScrollX(b)
	case addr == 0xFF44:
		g.resetLY()
	case addr == 0xFF47:
		g.bgp = b
	default:
		log.Panicf("unimplemented write gpu addr 0x%04X = 0x%02X\n", addr, b)
	}
}

func (g *GPU) ReadByte(addr uint16) byte {
	switch {
	case addr >= 0x8000 && addr <= 0x9FFF:
		return g.vram[addr-0x8000]
	case addr == 0xFF40:
		return g.getControl()
	case addr == 0xFF41:
		return g.getStat()
	case addr == 0xFF42:
		return g.getScrollY()
	case addr == 0xFF43:
		return g.getScrollX()
	case addr == 0xFF44:
		return g.getLY()
	case addr == 0xFF47:
		return g.bgp
	default:
		log.Panicf("unimplemented read gpu addr 0x%04X\n", addr)
	}
	panic("unexpected gpu failure")
}

type Dot byte

const (
	DotLighter Dot = iota
	DotLight
	DotDark
	DotDarker
)

func (g *GPU) getPalette(d Dot) byte {
	offset := 2 * d
	color := (g.bgp >> offset)
	return color & 0b11
}

func New() *GPU {
	return &GPU{
		vram: make([]byte, 8*1024),
	}
}

func (g *GPU) setScrollX(x byte) {
	g.scx = x
}

func (g *GPU) getScrollX() byte {
	return g.scx
}

func (g *GPU) setScrollY(y byte) {
	g.scy = y
}

func (g *GPU) getScrollY() byte {
	return g.scy
}

func (g *GPU) setStat(s byte) {
	g.stat = s
}

func (g *GPU) getStat() byte {
	return g.stat
}

func (g *GPU) setControl(b byte) {
	g.lcdEnable = b&(1<<7) > 0
	g.winTileMapArea = b&(1<<6) > 0
	g.winEnable = b&(1<<5) > 0
	g.bgAndWinTileDataArea = b&(1<<4) > 0
	g.bgTileMapArea = b&(1<<3) > 0
	g.objSize = b&(1<<2) > 0
	g.objEnable = b&(1<<1) > 0
	g.bgAndWinEnablePriority = b&1 > 0
}

func (g *GPU) getControl() byte {
	var b byte
	if g.lcdEnable {
		b |= 1 << 7
	}
	if g.winTileMapArea {
		b |= 1 << 6
	}
	if g.winEnable {
		b |= 1 << 5
	}
	if g.bgAndWinTileDataArea {
		b |= 1 << 4
	}
	if g.bgTileMapArea {
		b |= 1 << 3
	}
	if g.objSize {
		b |= 1 << 2
	}
	if g.objEnable {
		b |= 1 << 1
	}
	if g.bgAndWinEnablePriority {
		b |= 1 << 0
	}
	return b
}

func (g *GPU) resetLY() {
	g.ly = 0
}

func (g *GPU) getLY() byte {
	return 0x94
	// return g.ly
}

func (g *GPU) Run() {
	cfg := pixelgl.WindowConfig{
		Title:  "gameboy",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		g.render()
		win.Update()
	}
}

func (g *GPU) render() {
	g.renderBackground()
	g.renderWindow()
	g.renderSprites()
}

func (g *GPU) renderBackground() {
	// is bg turned on?

	// which tile map to use? 0x9800-0x9BFF or 0x9C00-0x9FF?

	// read tile map into 32x32 byte array
}

func (g *GPU) renderWindow() {

}

func (g *GPU) renderSprites() {

}
