package nes

// PC is not incremented in the opcode handler functions,
// to address the operand it must do the addition - cpu.rom[cpu.PC+1]
// For overflow protection integer type casts are used

// Opcode represents a single opcode
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
	Zpy                 // Zero page, Y
	Abs                 // Absolute
	Abx                 // Absolute, X
	Aby                 // Absolute, Y
	Ind                 // Indirect
	Izx                 // Indexed indirect, X
	Izy                 // Indirect indexed, Y
	Acc                 // Accumulator
	Imp                 // Implied
	Rel                 // Relative
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
	0x18: &Opcode{clc, Imp, 1, 2},
	0x58: &Opcode{cli, Imp, 1, 2},
	0xB8: &Opcode{clv, Imp, 1, 2},

	// STA
	0x85: &Opcode{sta, Zp, 2, 3},
	0x95: &Opcode{sta, Zpx, 2, 4},
	0x8D: &Opcode{sta, Abs, 3, 4},
	0x9D: &Opcode{sta, Abx, 3, 5},
	0x99: &Opcode{sta, Aby, 3, 5},
	0x81: &Opcode{sta, Izx, 2, 6},
	0x91: &Opcode{sta, Izy, 2, 6},

	0x86: &Opcode{stx, Zp, 2, 3},
	0x96: &Opcode{stx, Zpy, 2, 4},
	0x8E: &Opcode{stx, Abs, 3, 4},

	0x84: &Opcode{sty, Zp, 2, 3},
	0x94: &Opcode{sty, Zpx, 2, 4},
	0x8C: &Opcode{sty, Abs, 3, 4},

	// ADC
	0x69: &Opcode{adc, Imm, 2, 2},
	0x65: &Opcode{adc, Zp, 2, 3},
	0x75: &Opcode{adc, Zpx, 2, 4},
	0x6D: &Opcode{adc, Abs, 3, 4},
	0x7D: &Opcode{adc, Abx, 3, 4},
	0x79: &Opcode{adc, Aby, 3, 4},
	0x61: &Opcode{adc, Izx, 2, 6},
	0x71: &Opcode{adc, Izy, 2, 5},

	// AND
	0x29: &Opcode{and, Imm, 2, 2},
	0x25: &Opcode{and, Zp, 2, 3},
	0x35: &Opcode{and, Zpx, 2, 4},
	0x2D: &Opcode{and, Abs, 3, 4},
	0x3D: &Opcode{and, Abx, 3, 4},
	0x39: &Opcode{and, Aby, 3, 4},
	0x21: &Opcode{and, Izx, 2, 6},
	0x31: &Opcode{and, Izy, 2, 5},

	// ASL
	0x0A: &Opcode{asl, Acc, 1, 2},
	0x06: &Opcode{asl, Zp, 2, 5},
	0x16: &Opcode{asl, Zpx, 2, 6},
	0x0E: &Opcode{asl, Abs, 3, 6},
	0x1E: &Opcode{asl, Abx, 3, 7},

	// Branch instructions
	0x90: &Opcode{bcc, Rel, 2, 2},
	0xB0: &Opcode{bcs, Rel, 2, 2},
	0xF0: &Opcode{beq, Rel, 2, 2},
	0x30: &Opcode{bmi, Rel, 2, 2},
	0xD0: &Opcode{bne, Rel, 2, 2},
	0x10: &Opcode{bpl, Rel, 2, 2},
	0x50: &Opcode{bvc, Rel, 2, 2},
	0x70: &Opcode{bvs, Rel, 2, 2},

	0x00: &Opcode{brk, Imp, 1, 7},

	0x24: &Opcode{bit, Zp, 2, 3},
	0x2C: &Opcode{bit, Abs, 3, 4},

	// Comparison instructions
	0xC9: &Opcode{cmp, Imm, 2, 2},
	0xC5: &Opcode{cmp, Zp, 2, 3},
	0xD5: &Opcode{cmp, Zpx, 2, 4},
	0xCD: &Opcode{cmp, Abs, 3, 4},
	0xDD: &Opcode{cmp, Abx, 3, 4},
	0xD9: &Opcode{cmp, Aby, 3, 4},
	0xC1: &Opcode{cmp, Izx, 2, 6},
	0xD1: &Opcode{cmp, Izy, 2, 5},

	0xE0: &Opcode{cpx, Imm, 2, 2},
	0xE4: &Opcode{cpx, Zp, 2, 3},
	0xEC: &Opcode{cpx, Abs, 3, 4},

	0xC0: &Opcode{cpy, Imm, 2, 2},
	0xC4: &Opcode{cpy, Zp, 2, 3},
	0xCC: &Opcode{cpy, Abs, 3, 4},

	// Increment / decrement instructions
	0xC6: &Opcode{dec, Zp, 2, 5},
	0xD6: &Opcode{dec, Zpx, 2, 6},
	0xCE: &Opcode{dec, Abs, 3, 6},
	0xDE: &Opcode{dec, Abx, 3, 7},

	0xE6: &Opcode{inc, Zp, 2, 5},
	0xF6: &Opcode{inc, Zpx, 2, 6},
	0xEE: &Opcode{inc, Abs, 3, 6},
	0xFE: &Opcode{inc, Abx, 3, 7},

	0xCA: &Opcode{dex, Imp, 1, 2},
	0x88: &Opcode{dey, Imp, 1, 2},

	0xE8: &Opcode{inx, Imp, 1, 2},
	0xC8: &Opcode{iny, Imp, 1, 2},

	// EOR
	0x49: &Opcode{eor, Imm, 2, 2},
	0x45: &Opcode{eor, Zp, 2, 3},
	0x55: &Opcode{eor, Zpx, 2, 4},
	0x4D: &Opcode{eor, Abs, 3, 4},
	0x5D: &Opcode{eor, Abx, 3, 4},
	0x59: &Opcode{eor, Aby, 3, 4},
	0x41: &Opcode{eor, Izx, 2, 6},
	0x51: &Opcode{eor, Izy, 2, 5},

	0x20: &Opcode{jsr, Abs, 3, 6},

	0xA2: &Opcode{ldx, Imm, 2, 2},
	0xA6: &Opcode{ldx, Zp, 2, 3},
	0xB6: &Opcode{ldx, Zpy, 2, 4},
	0xAE: &Opcode{ldx, Abs, 3, 4},
	0xBE: &Opcode{ldx, Aby, 3, 4},

	0xA0: &Opcode{ldy, Imm, 2, 2},
	0xA4: &Opcode{ldy, Zp, 2, 3},
	0xB4: &Opcode{ldy, Zpx, 2, 4},
	0xAC: &Opcode{ldy, Abs, 3, 4},
	0xBC: &Opcode{ldy, Abx, 3, 4},

	0x4A: &Opcode{lsr, Acc, 1, 2},
	0x46: &Opcode{lsr, Zp, 2, 5},
	0x56: &Opcode{lsr, Zpx, 2, 6},
	0x4E: &Opcode{lsr, Abs, 3, 6},
	0x5E: &Opcode{lsr, Abx, 3, 7},

	0xEA: &Opcode{nop, Imp, 1, 2},

	0x09: &Opcode{ora, Imm, 2, 2},
	0x05: &Opcode{ora, Zp, 2, 3},
	0x15: &Opcode{ora, Zpx, 2, 4},
	0x0D: &Opcode{ora, Abs, 3, 4},
	0x1D: &Opcode{ora, Abx, 3, 4},
	0x19: &Opcode{ora, Aby, 3, 4},
	0x01: &Opcode{ora, Izx, 2, 6},
	0x11: &Opcode{ora, Izy, 2, 5},

	0x48: &Opcode{pha, Imp, 1, 3},
	0x08: &Opcode{php, Imp, 1, 3},
	0x68: &Opcode{pla, Imp, 1, 4},
	0x28: &Opcode{plp, Imp, 1, 4},

	0x2A: &Opcode{rol, Acc, 1, 2},
	0x26: &Opcode{rol, Zp, 2, 5},
	0x36: &Opcode{rol, Zpx, 2, 6},
	0x2E: &Opcode{rol, Abs, 3, 6},
	0x3E: &Opcode{rol, Abx, 3, 7},

	0x6A: &Opcode{ror, Acc, 1, 2},
	0x66: &Opcode{ror, Zp, 2, 5},
	0x76: &Opcode{ror, Zpx, 2, 6},
	0x6E: &Opcode{ror, Abs, 3, 6},
	0x7E: &Opcode{ror, Abx, 3, 7},

	0x40: &Opcode{rti, Imp, 1, 6},
	0x60: &Opcode{rts, Imp, 1, 6},

	// SBC
	0xE9: &Opcode{sbc, Imm, 2, 2},
	0xE5: &Opcode{sbc, Zp, 2, 3},
	0xF5: &Opcode{sbc, Zpx, 2, 4},
	0xED: &Opcode{sbc, Abs, 3, 4},
	0xFD: &Opcode{sbc, Abx, 3, 4},
	0xF9: &Opcode{sbc, Aby, 3, 4},
	0xE1: &Opcode{sbc, Izx, 2, 6},
	0xF1: &Opcode{sbc, Izy, 2, 5},

	0x38: &Opcode{sec, Imp, 1, 2},
	0xF8: &Opcode{sed, Imp, 1, 2},

	0xAA: &Opcode{tax, Imp, 1, 2},
	0xA8: &Opcode{tay, Imp, 1, 2},
	0xBA: &Opcode{tsx, Imp, 1, 2},
	0x8A: &Opcode{txa, Imp, 1, 2},
	0x9A: &Opcode{txs, Imp, 1, 2},
	0x98: &Opcode{tya, Imp, 1, 2},
}

// LDA (Load Accumulator)
// Loads a byte into the accumulator setting the zero and negative
// flags as appropriate
func lda(c *CPU, m byte) {
	c.A = c.mem.Read(peek(c, m))

	if c.A>>7 == 1 {
		c.P |= 0x80
	} else if c.A == 0 {
		c.P |= 0x2
	}
}

func jmp(c *CPU, m byte) {
	if m == Ind {
		a := (uint16(c.mem.Read(c.PC+2)) << 8) | uint16(c.mem.Read(c.PC+1))
		c.PC = (uint16(c.mem.Read(a+1)) << 8) | uint16(c.mem.Read(a))
	} else {
		c.PC = (uint16(c.mem.Read(c.PC+2)) << 8) | uint16(c.mem.Read(c.PC+1))
	}
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

func adc(c *CPU, m byte) {
	v := c.mem.Read(peek(c, m))
	a := c.A
	c.A = a + v + (c.P & 0x01)
	c.testOverflowOnAdd(a, v, c.A)
	c.testNegative(c.A)
	c.testZero(c.A)
	c.testCarryOnAdd(c.A)
}

func and(c *CPU, m byte) {
	c.A &= c.mem.Read(peek(c, m))
	c.testNegative(c.A)
	c.testZero(c.A)
}

func asl(c *CPU, m byte) {
	if m == Acc {
		if c.A&0x80 > 0 {
			c.setFlag(FlagCarry)
		} else {
			c.clearFlag(FlagCarry)
		}
		c.A = c.A << 1
		c.testZero(c.A)
		c.testNegative(c.A)
	} else {
		addr := peek(c, m)
		v := c.mem.Read(addr)
		if v&0x80 > 0 {
			c.setFlag(FlagCarry)
		} else {
			c.clearFlag(FlagCarry)
		}
		v = v << 1
		c.mem.Write(addr, v)
		c.testZero(v)
		c.testNegative(v)
	}
}

func bcc(c *CPU, m byte) {
	// FIXME: implement checks
	c.PC += uint16(peek(c, m))
}

func bcs(c *CPU, m byte) {

}

func beq(c *CPU, m byte) {

}

func bmi(c *CPU, m byte) {

}

func bne(c *CPU, m byte) {

}

func bpl(c *CPU, m byte) {

}

func bvc(c *CPU, m byte) {

}

func bvs(c *CPU, m byte) {

}

func bit(c *CPU, m byte) {

}

func brk(c *CPU, m byte) {
}

func clc(c *CPU, m byte) {
}

func cli(c *CPU, m byte) {
}

func clv(c *CPU, m byte) {
}

func cmp(c *CPU, m byte) {

}

func cpx(c *CPU, m byte) {

}

func cpy(c *CPU, m byte) {

}

func inc(c *CPU, m byte) {

}

func inx(c *CPU, m byte) {

}

func iny(c *CPU, m byte) {

}

func dec(c *CPU, m byte) {

}

func dex(c *CPU, m byte) {

}

func dey(c *CPU, m byte) {

}

func eor(c *CPU, m byte) {

}

func jsr(c *CPU, m byte) {
}

func ldx(c *CPU, m byte) {
}

func ldy(c *CPU, m byte) {
}

func lsr(c *CPU, m byte) {
}

func rol(c *CPU, m byte) {
}

func ror(c *CPU, m byte) {
}

func nop(c *CPU, m byte) {
}

func ora(c *CPU, m byte) {

}

func pha(c *CPU, m byte) {

}

func php(c *CPU, m byte) {

}

func pla(c *CPU, m byte) {

}

func plp(c *CPU, m byte) {

}

func rti(c *CPU, m byte) {

}

func rts(c *CPU, m byte) {

}

func sbc(c *CPU, m byte) {

}

func sec(c *CPU, m byte) {

}

func sed(c *CPU, m byte) {

}

func stx(c *CPU, m byte) {

}

func sty(c *CPU, m byte) {

}

func tax(c *CPU, m byte) {

}

func tay(c *CPU, m byte) {

}

func tsx(c *CPU, m byte) {

}

func txa(c *CPU, m byte) {

}

func txs(c *CPU, m byte) {

}

func tya(c *CPU, m byte) {

}
