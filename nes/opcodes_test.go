package nes

import (
	"testing"
)

func TestSei(t *testing.T) {
	m := Memory{mapper: &NROM128{}}

	// Initialize memory and CPU
	var prgBegin uint16 = 0x8000
	m.ram[prgBegin] = 0x78
	m.ram[prgBegin+1] = 0x58

	m.Write(0xFFFD, byte(prgBegin>>8))
	m.Write(0xFFFC, byte(prgBegin&0xFF))

	cpu := InitCPU(&m)
	cpu.Reset()
	cpu.clearFlag(FlagInterruptDisable)

	// Call SEI instruction handler
	sei(cpu, Imp)

	// Assert
	if cpu.P&FlagInterruptDisable != FlagInterruptDisable {
		t.Fatalf("SEI failed to set I flag %08b", cpu.P)
	}
}

func TestCli(t *testing.T) {
	m := Memory{mapper: &NROM128{}}

	// Initialize memory and CPU
	var prgBegin uint16 = 0x8000
	m.ram[prgBegin] = 0x78
	m.ram[prgBegin+1] = 0x58

	m.Write(0xFFFD, byte(prgBegin>>8))
	m.Write(0xFFFC, byte(prgBegin&0xFF))

	cpu := InitCPU(&m)
	cpu.Reset()

	// Call CLI instruction handler
	cli(cpu, Imp)

	// Assert
	if cpu.P&FlagInterruptDisable == FlagInterruptDisable {
		t.Fatalf("CLI failed to clear I flag %08b", cpu.P)
	}
}
