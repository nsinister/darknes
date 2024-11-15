// Package nes implements NES emulation
package nes

import (
	"fmt"
	"reflect"
	"runtime"
)

const (
	FlagCarry            byte = 0x01
	FlagZero             byte = 0x02
	FlagInterruptDisable byte = 0x04
	FlagDecimalMode      byte = 0x08
	FlagBreakCommand     byte = 0x10
	FlagOverflow         byte = 0x40
	FlagNegative         byte = 0x80

	prgStartAddr = 0x8000

	stackAddr uint16 = 0x0100
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

	cyclesPassed uint64

	// RAM
	mem *Memory
}

type CpuState struct {
	A            byte
	X, Y         byte
	S            byte
	PC           uint16
	P            byte
	Cycles       uint16
	CyclesPassed uint64
	// Last instruction
	LastOp *Opcode
}

// InitCPU mehtod initializes 2A03 CPU.
// Returns CPU struct with allocated memory.
func InitCPU(m *Memory) *CPU {
	return &CPU{mem: m}
}

// Reset method resets the CPU to match its power up state.
func (cpu *CPU) Reset() {
	//cpu.cyclesPassed = 0
	//cpu.cycles = 0

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
	cpu.PC = a
}

// Step executes a single instruction
func (cpu *CPU) Step() CpuState {
	// Read next instruction
	op := cpu.mem.Read(cpu.PC)

	// Identify and process the instruction
	opcode, ok := opcodeMap[op]
	if ok {
		fmt.Printf("OP=%x Handler: %s\n", op, getFunctionName(opcode.Handler))

		// Execute instruction
		opcode.Handler(cpu, opcode.mode)
	} else {
		panic(fmt.Sprintf("Opcode %x not recognized\n", op))
	}

	cpu.cycles += opcode.cycles
	cpu.cyclesPassed += uint64(cpu.cycles)

	// Return current copy of CPU state for debugging
	state := CpuState{A: cpu.A,
		X:            cpu.X,
		Y:            cpu.Y,
		PC:           cpu.PC,
		P:            cpu.P,
		S:            cpu.S,
		Cycles:       cpu.cycles,
		CyclesPassed: cpu.cyclesPassed,
		LastOp:       opcode}

	// Go to the next instruction (if prev opcode was jmp, it doesn't do anything)
	nextOp(cpu, op)

	return state
}

func (cpu *CPU) nmiInterrupt() {
	cpu.pushWord(cpu.PC)
	cpu.push(cpu.P)
	// TODO unfinished
	cpu.setFlag(FlagInterruptDisable)

	// fetch address vector
	pcLow := cpu.mem.Read(0xFFFA)
	pcHigh := cpu.mem.Read(0xFFFB)
	// jump to the address
	cpu.PC = uint16((pcHigh << 8) | pcLow)
}

// Push value to stack
func (cpu *CPU) push(val byte) {
	addr := stackAddr + uint16(cpu.S)
	cpu.mem.Write(addr, val)
	cpu.S--
}

// Pop value from stack
func (cpu *CPU) pop() byte {
	addr := stackAddr + uint16(cpu.S)
	cpu.S++
	return cpu.mem.Read(addr)
}

func (cpu *CPU) pushWord(val uint16) {
	cpu.push(byte(val >> 8))
	cpu.push(byte(val & 0xFF))
}

func (cpu *CPU) popWord() uint16 {
	lowByte := cpu.pop()
	highByte := cpu.pop()
	return uint16(highByte<<8 | lowByte)
}

// Sets flag in register P
func (cpu *CPU) setFlag(flag byte) {
	cpu.P |= flag
}

func (cpu *CPU) clearFlag(flag byte) {
	cpu.P &= ^flag
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

// Sets PC (Program Counter) to the beginning of the next instruction
func nextOp(cpu *CPU, opcode byte) {
	// Jump, branch and similar instructions skip the PC increment
	switch opcode {
	case 0x00:
		fallthrough
	case 0x40:
		fallthrough
	case 0x60:
		fallthrough
	case 0x20:
		fallthrough
	case 0x4C:
		fallthrough
	case 0x6C:
		fallthrough
	case 0x90:
		fallthrough
	case 0xB0:
		fallthrough
	case 0xF0:
		fallthrough
	case 0x30:
		fallthrough
	case 0xD0:
		fallthrough
	case 0x10:
		fallthrough
	case 0x50:
		fallthrough
	case 0x70:
		return
	}
	cpu.PC += opcodeMap[opcode].length
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
