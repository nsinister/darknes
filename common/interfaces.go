package common

import "darknes/nes"

type CpuEmulator interface {
	Step() nes.CpuState
}

type PpuEmulator interface {
	Step(cycles uint16)
}
