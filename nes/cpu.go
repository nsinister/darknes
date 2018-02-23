// Package nes implements NES emulation
package nes

import (
	"fmt"
	"reflect"
	"runtime"
)

const (
	FlagCarry    byte = 0x01
	FlagZero     byte = 0x02
	FlagOverflow byte = 0x40
	FlagNegative byte = 0x80

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

	// Cycles to wait
	cycles uint16

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
	a := (uint16(cpu.mem.Read(0xFFFD)) << 8) | uint16(cpu.mem.Read(0xFFFC))
	cpu.PC = (uint16(cpu.mem.Read(a+1)) << 8) | uint16(cpu.mem.Read(a))
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
		fmt.Printf("OP=%x Handler: %s\n", op, getFunctionName(opcode.handler))
		opcode.handler(cpu, opcode.mode)
	} else {
		panic(fmt.Sprintf("Opcode %x not recognized\n", op))
	}
	// Go to the next instruction (if prev opcode was jmp, it doesn't do anything)
	nextOp(cpu, op)
}

// Sets flag in register P
func (cpu *CPU) setFlag(flag byte) {
	cpu.P |= flag
}

func (cpu *CPU) clearFlag(flag byte) {
	switch flag {
	case FlagOverflow:
		cpu.P &= 0xBF
	case FlagNegative:
		cpu.P &= 0x7F
	case FlagZero:
		cpu.P &= 0xFD
	}
}

func (cpu *CPU) setBranchCycles(addr uint16) {
	if ((cpu.PC + 1) & 0xFF00 >> 8) != ((addr & 0xFF00) >> 8) {
		cpu.cycles = 4
	} else {
		cpu.cycles = 3
	}
}

func (cpu *CPU) testOverflowOnAdd(val1 byte, val2 byte, res byte) {
	if ((val1^val2)&0x80 == 0x0) && ((val1^res)&0x80 == 0x80) {
		cpu.setFlag(FlagOverflow)
		return
	}
	cpu.clearFlag(FlagOverflow)
}

func (cpu *CPU) testOverflowOnSub(val1 byte, val2 byte) {
	r := val1 - val2 - (1 - cpu.P&0x01)
	if ((val1^r)&0x80) != 0 && ((val1^val2)&0x80) != 0 {
		cpu.setFlag(FlagOverflow)
		return
	}
	cpu.clearFlag(FlagOverflow)
}

func (cpu *CPU) testNegative(val byte) {
	if val&0x80 == 0x80 {
		cpu.setFlag(FlagNegative)
		return
	}
	cpu.clearFlag(FlagNegative)
}

func (cpu *CPU) testZero(val byte) {
	if val == 0 {
		cpu.setFlag(FlagZero)
		return
	}
	cpu.clearFlag(FlagZero)
}

func (cpu *CPU) testCarryOnAdd(val byte) {
	if val > 0xFF {
		cpu.setFlag(FlagCarry)
		return
	}
	cpu.clearFlag(FlagCarry)
}

// Sets PC (Program Counter) to the beginning of the next instruction.
// Depending on the instruction, it may jump to a specific address
// or increment PC by instruction length.
func nextOp(cpu *CPU, opcode byte) {
	// JMP or equivalent instructions skip the PC increment
	if opcode == 0x4C || opcode == 0x6C {
		return
	}
	cpu.PC += opcodeMap[opcode].len
}

// Takes address value by specified addressing mode using instruction operand
func peek(c *CPU, m byte) uint16 {
	var retval uint16
	switch m {
	case Imm:
		retval = c.PC + 1
	case Zp:
		retval = uint16(c.mem.Read(c.PC + 1))
	case Zpx:
		retval = uint16(byte((c.mem.Read(c.PC+1) + c.X)))
	case Zpy:
		retval = uint16(byte((c.mem.Read(c.PC+1) + c.Y)))
	case Abs:
		retval = uint16((uint16(c.mem.Read(c.PC+2)) << 8) | uint16(c.mem.Read(c.PC+1)))
	case Abx:
		addr := (uint16(c.mem.Read(c.PC+2)) << 8) | uint16(c.mem.Read(c.PC+1))
		// +1 if page crossed
		if addr&0xFF00 != (addr+uint16(c.X))&0xFF00 {
			c.cycles++
		}
		retval = uint16(addr + uint16(c.X))
	case Aby:
		addr := (uint16(c.mem.Read(c.PC+2)) << 8) | uint16(c.mem.Read(c.PC+1))
		// +1 if page crossed
		if addr&0xFF00 != (addr+uint16(c.Y))&0xFF00 {
			c.cycles++
		}
		retval = uint16(addr + uint16(c.Y))
	case Izx:
		a := uint16(byte((c.mem.Read(c.PC+1) + c.X)))
		retval = uint16((uint16(c.mem.Read(a+1)) << 8) | uint16(c.mem.Read(a)))
	case Izy:
		a := uint16(c.mem.Read(c.PC + 1))
		addr := (uint16(c.mem.Read(a+1)) << 8) | uint16(c.mem.Read(a))
		// +1 if page crossed
		if addr&0xFF00 != (addr+uint16(c.Y))&0xFF00 {
			c.cycles++
		}
		retval = addr + uint16(c.Y)
	// Indirect is used only by JMP, so the below branch will be left unused
	case Ind:
		a := (uint16(c.mem.Read(c.PC+2)) << 8) | uint16(c.mem.Read(c.PC+1))
		retval = (uint16(c.mem.Read(a+1)) << 8) | uint16(c.mem.Read(a))

	//case Acc:
	//	retval = c.A
	case Rel:
		retval = c.PC + 1
	default:
		panic(fmt.Sprintf("Addressing mode %v not recognized\n", m))
	}
	return retval
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
