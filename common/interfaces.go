package common

import "darknes/nes"

type CpuEmulator interface {
	GetCpu() *nes.CPU
	Step()
}
