package main

import (
	"flag"
	"log"

	_ "embed"

	"github.com/faiface/pixel/pixelgl"
	"github.com/prestonp/gbc/pkg/gb"
	"github.com/prestonp/gbc/pkg/gb/apu"
	"github.com/prestonp/gbc/pkg/gb/gpu"
)

//go:embed boot.gb
var boot []byte

var (
	debug = flag.Bool("debug", false, "debug mode")
	file  = flag.String("f", "", "rom file")
)

func main() {
	flag.Parse()

	if *file == "" {
		log.Fatal("missing filename")
	}
	rom, err := gb.ReadRom(*file)
	if err != nil {
		log.Fatal(err)
	}

	gpu := gpu.New()
	apu := apu.New()
	mmu := gb.NewMMU(boot, rom, gpu, apu)
	cpu := gb.NewCPU(mmu, gpu, *debug)

	pixelgl.Run(cpu.Run)
}
