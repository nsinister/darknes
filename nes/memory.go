package nes

import "fmt"

// Mapper represents memory mapping
type Mapper interface {
	Translate(uint16) uint16
}

// Memory represents NES memory model
type Memory struct {
	// RAM
	ram    [65536]byte
	mapper Mapper
}

// GetMapper returns iNES mapper
func GetMapper(h *RomHeader) Mapper {
	switch h.MapperNum {
	case 0:
		if h.PrgRomSize == 1 {
			return &NROM128{}
		}
		return &NROM256{}
	default:
		panic(fmt.Sprintf("Mapper %v not implemented\n", h.MapperNum))
	}
}

// Load loads NES ROM into NES memory and returns Memory
func (rom *Rom) Load() *Memory {
	mp := GetMapper(rom.Header)
	m := Memory{mapper: mp}
	var p = 0x8000
	for i := range rom.prgRom {
		m.ram[p] = rom.prgRom[i]
		p++
	}

	// Write reset vector into 0xFFFC
	m.Write(0xFFFD, byte(0x8000>>8))
	m.Write(0xFFFC, byte(0x8000&0xFF))

	return &m
}

// Translate performs mirroring and mapping of the address where needed
// and returns effective address
func (m *Memory) Translate(addr uint16) uint16 {
	switch {
	// NES RAM mirroring
	case addr > 0x07FF && addr <= 0x1FFF:
		return addr & 0x07FF
	// NES PPU registers mirroring
	case addr > 0x2007 && addr <= 0x3FFF:
		return addr & 0x2007
	default:
		return m.mapper.Translate(addr)
	}
}

func (m *Memory) Read(addr uint16) byte {
	return m.ram[m.Translate(addr)]
}

func (m *Memory) Write(addr uint16, val byte) {
	m.ram[m.Translate(addr)] = val
}
