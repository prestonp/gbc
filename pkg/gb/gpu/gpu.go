package gpu

import (
	"fmt"
	"log"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type GPU struct {
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
	fmt.Fprintf(&b, "\tlcd control: %08b\n", g.GetControl())
	return b.String()
}

func (g *GPU) SetRegister(addr uint16, b byte) {
	switch {
	case addr == 0xFF47:
		g.bgp = b
	default:
		log.Panicf("unimplemented gpu register 0x%04X = 0x%02X\n", addr, b)
	}
}

func (g *GPU) GetRegister(addr uint16) byte {
	switch {
	case addr == 0xFF47:
		return g.bgp
	default:
		log.Panicf("unimplemented gpu register 0x%04X\n", addr)
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
	return &GPU{}
}

func (g *GPU) SetScrollX(x byte) {
	g.scx = x
}

func (g *GPU) GetScrollX() byte {
	return g.scx
}

func (g *GPU) SetScrollY(y byte) {
	g.scy = y
}

func (g *GPU) GetScrollY() byte {
	return g.scy
}

func (g *GPU) SetStat(s byte) {
	g.stat = s
}

func (g *GPU) GetStat() byte {
	return g.stat
}

func (g *GPU) SetControl(b byte) {
	g.lcdEnable = b&(1<<7) > 0
	g.winTileMapArea = b&(1<<6) > 0
	g.winEnable = b&(1<<5) > 0
	g.bgAndWinTileDataArea = b&(1<<4) > 0
	g.bgTileMapArea = b&(1<<3) > 0
	g.objSize = b&(1<<2) > 0
	g.objEnable = b&(1<<1) > 0
	g.bgAndWinEnablePriority = b&1 > 0
}

func (g *GPU) GetControl() byte {
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

func (g *GPU) ResetLY() {
	g.ly = 0
}

func (g *GPU) GetLY() byte {
	return 0x94
	// return g.ly
}

func (g *GPU) Loop(update func()) {
	cfg := pixelgl.WindowConfig{
		Title:  "gameboy",
		Bounds: pixel.R(0, 0, 1024, 768),
		// VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		update()
		g.render()
		win.Update()
	}
}

func (g *GPU) render() {

}
