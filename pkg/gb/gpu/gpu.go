package gpu

import (
	"fmt"
	"image"
	"image/color"
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
		g.render(win)
		win.Update()
	}
}

func (g *GPU) render(win *pixelgl.Window) {
	g.renderBackground(win)
	g.renderWindow(win)
	g.renderSprites(win)
}

func (g *GPU) renderBackground(win *pixelgl.Window) {
	if !g.lcdEnable {
		return
	}

	if !g.bgAndWinEnablePriority {
		win.Clear(color.White)
		return
	}

	pic := pixel.PictureDataFromImage(g)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
}

func (g *GPU) renderWindow(win *pixelgl.Window) {

}

func (g *GPU) renderSprites(win *pixelgl.Window) {

}

// read a tile into a byte slice storing the color IDs. The color IDs must
// refer to palette to produce actual colors. The slice is flat, but is indexed
// in row, col order.
func (g *GPU) readTile(addrMode uint16, idx byte) []byte {
	if addrMode == 0x8800 {
		log.Panicf("0x8800 addr mode not implemented")
	}

	baseAddr := addrMode + uint16(idx)*16

	var b []byte

	// read 16 bytes, each pair represents a line, refer to gb spec/docs for encoding
	for row := uint16(0); row < 8; row++ {
		lower := g.ReadByte(baseAddr + (row * 2))
		upper := g.ReadByte(baseAddr + (row*2 + 1))

		for col := 0; col < 8; col++ {
			offset := 7 - col
			mask := byte(1 << offset)
			colorId := upper&mask>>(offset+1) | (lower&mask)>>offset
			b = append(b, colorId)
		}
	}
	return b
}

var _ image.Image = &GPU{}

func (g *GPU) At(x, y int) color.Color {
	tileR := y / 8
	tileC := x / 8

	tileIdx := tileR*32 + tileC

	addrMode := func() uint16 {
		if g.bgAndWinTileDataArea {
			return 0x8000
		}
		return 0x8800
	}()

	tileMapOffset := uint16(0x9800)
	if g.bgTileMapArea {
		tileMapOffset = 0x9C00
	}
	tileAddr := tileMapOffset + uint16(tileIdx)
	tileID := g.ReadByte(tileAddr)
	tileData := g.readTile(addrMode, tileID)

	tileX := x % 8
	tileY := y % 8

	colorID := tileData[tileY*8+tileX]

	// todo: use the color palette
	switch colorID {
	case 0:
		return color.RGBA{255, 255, 255, 255}
	case 1:
		return color.RGBA{150, 150, 150, 255}
	case 2:
		return color.RGBA{50, 50, 50, 255}
	case 3:
		return color.RGBA{0, 0, 0, 255}
	default:
		panic("unknown color id")
	}
}

func (g *GPU) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{255, 255},
	}
}

func (g *GPU) ColorModel() color.Model {
	return color.RGBAModel
}
