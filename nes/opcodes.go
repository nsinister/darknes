package nes

import (
	"fmt"
)

// PC is not incremented in the opcode handler functions,
// to address the operand it must do the addition - cpu.rom[cpu.PC+1]
// For overflow protection integer type casts are used

// Represents a single opcode
type Opcode struct {
	// Function that handles the instruction
	handler func(*CPU, byte)
	// Addressing mode
	mode byte
	// Instruction length
	len uint16
	// Number of cycles
	cycles uint16
}

const (
	// Addressing modes
	Imm byte = iota + 1 // Immediate
	Zp                  // Zero page
	Zpx                 // Zero page, X
	Abs                 // Absolute
	Abx                 // Absolute, X
	Aby                 // Absolute, Y
	Ind                 // Indirect
	Izx                 // Indexed indirect, X
	Izy                 // Indirect indexed, Y
	Acc
	Imp

	// Masks for overflow control
	u8mask  = 0x00FF
	u16mask = 0x0000FFFF
)

// Opcode map contains bindings between opcode functions and their
// hexdecimal representation

var opcodeMap = map[byte]*Opcode{
	0xA9: &Opcode{lda, Imm, 2, 2},
	0xA5: &Opcode{lda, Zp, 2, 3},
	0xB5: &Opcode{lda, Zpx, 2, 4},
	0xAD: &Opcode{lda, Abs, 3, 4},
	0xBD: &Opcode{lda, Abx, 3, 4},
	0xB9: &Opcode{lda, Aby, 3, 4},
	0xA1: &Opcode{lda, Izx, 2, 6},
	0xB1: &Opcode{lda, Izy, 2, 5},

	0x4C: &Opcode{jmp, Abs, 3, 3},
	0x6C: &Opcode{jmp, Ind, 3, 5},

	0x78: &Opcode{sei, Imp, 1, 2},
	0xD8: &Opcode{cld, Imp, 1, 2},

	// STA
	0x8D: &Opcode{sta, Abs, 3, 4},
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

// Takes one byte from RAM using specified addressing mode
func peek(c *CPU, m byte) byte {
	// TODO: Review this, write tests
	switch m {
	case Imm:
		return c.mem.Read(c.PC + 1)
	case Zp:
		return c.mem.Read(uint16(c.mem.Read(c.PC + 1)))
	case Zpx:
		return c.mem.Read(uint16(byte((c.mem.Read(c.PC+1) + c.X) & u8mask)))
	case Abs:
		return c.mem.Read(uint16((uint16(c.mem.Read(c.PC+2)) << 8) | uint16(c.mem.Read(c.PC+1))))
	case Abx:
		return c.mem.Read(uint16(((uint16(c.mem.Read(c.PC+2))<<8)|uint16(c.mem.Read(c.PC+1)))+uint16(c.X)) & u16mask)
	case Aby:
		return c.mem.Read(uint16(((uint16(c.mem.Read(c.PC+2))<<8)|uint16(c.mem.Read(c.PC+1)))+uint16(c.Y)) & u16mask)
	case Izx:
		a := uint16(byte((c.mem.Read(c.PC+1) + c.X) & u8mask))
		addr := (uint16(c.mem.Read(a+1)) << 8) | uint16(c.mem.Read(a))
		return c.mem.Read(uint16(addr) & u16mask)
	case Izy:
		a := (uint16(c.mem.Read(c.PC+2)) << 8) | uint16(c.mem.Read(c.PC+1))
		return c.mem.Read(uint16((a + uint16(c.Y))) & u16mask)
	case Ind:
		a := (uint16(c.mem.Read(c.PC+2)) << 8) | uint16(c.mem.Read(c.PC+1))
		return c.mem.Read(a & u16mask)
	}
	panic(fmt.Sprintf("Addressing mode %v not recognized\n", m))
}

// LDA (Load Accumulator)
// Loads a byte into the accumulator setting the zero and negative
// flags as appropriate
func lda(c *CPU, m byte) {
	c.A = peek(c, m)

	if c.A>>7 == 1 {
		c.P |= 0x80
	} else if c.A == 0 {
		c.P |= 0x2
	}
}

func jmp(c *CPU, m byte) {
	c.PC = (uint16(c.mem.Read(c.PC+2)) << 8) | uint16(c.mem.Read(c.PC+1))
}

func sei(c *CPU, m byte) {
	// TODO
}

func cld(c *CPU, m byte) {
	// TODO
}

func sta(c *CPU, m byte) {
	// TODO
}
