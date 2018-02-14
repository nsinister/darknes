// Package nes implements NES emulation
package nes

import (
	"fmt"
)

const (
	prgStartAddr = 0x8000
)

// CPU represents 2A03 CPU based on MOS 6502
type CPU struct {
	// Accumulator register
	A byte
	// Index registers
	X, Y byte
	// Stack pointer register
	S byte
	// Program counter
	PC uint16
	// Status register
	P byte

	// RAM
	mem *Memory
}

// InitCPU mehtod initializes 2A03 CPU.
// Returns CPU struct with allocated memory.
func InitCPU(m *Memory) *CPU {
	return &CPU{mem: m}
}

// Reset method resets the CPU to match its power up state.
func (cpu *CPU) Reset() {
	cpu.X, cpu.Y = 0, 0
	cpu.P = 0x34
	cpu.S = 0xFD
	cpu.mem.Write(0x4017, 0x00)
	cpu.mem.Write(0x4015, 0x00)
	for i := uint16(0x4000); i <= 0x400F; i++ {
		cpu.mem.Write(i, 0x00)
	}

	// JMP (FFFC) - reset vector
	vec := (uint16(cpu.mem.Read(0xFFFD)) << 8) | uint16(cpu.mem.Read(0xFFFC))
	cpu.PC = vec
	/*
		cpu.S -= 3
		cpu.mem.Write(0x4015, 0x00)
		cpu.P |= 0x04
		cpu.PC = prgStartAddr */
}

// Step executes a single instruction
func (cpu *CPU) Step() {
	// Read next instruction
	op := cpu.mem.Read(cpu.PC)
	// Identify and process the instruction
	if opcode, ok := opcodeMap[op]; ok {
		opcode.handler(cpu, opcode.mode)
	} else {
		panic(fmt.Sprintf("Opcode %x not recognized\n", op))
	}
	// Go to the next instruction (if prev opcode was jmp, it doesn't do anything)
	nextOp(cpu, op)
}
