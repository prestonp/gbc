package main

import (
	"flag"
	"log"

	"github.com/prestonp/gbc/pkg/gb"
	"github.com/prestonp/gbc/pkg/gb/gpu"
)

var (
	debug = flag.Bool("debug", false, "debug mode")
	file  = flag.String("f", "", "rom file")
	boot  = flag.String("b", "bin/boot.gb", "path to boot rom")
)

func main() {
	flag.Parse()

	if *file == "" {
		log.Fatal("missing filename")
	}
	cartRom, err := gb.ReadRom(*file)
	if err != nil {
		log.Fatal(err)
	}

	if *boot == "" {
		log.Fatal("missing boot rom")
	}

	bootRom, err := gb.ReadRom(*boot)
	if err != nil {
		log.Fatal(err)
	}

	gpu := gpu.New()
	mmu := gb.NewMMU(bootRom, cartRom, gpu)
	cpu := gb.NewCPU(mmu, *debug)

	cpu.Run()
}
