package gb

import (
	"github.com/prestonp/gbc/pkg/gb/apu"
	"github.com/prestonp/gbc/pkg/gb/gpu"
)

func getTestCPU() *CPU {
	gpu := gpu.New()
	apu := apu.New()
	mmu := NewMMU(nil, nil, gpu, apu)
	return NewCPU(mmu, gpu, false)
}
