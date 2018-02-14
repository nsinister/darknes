package nes

// NROM128 is iNES mapper 000 with 1 PRG ROM bank
type NROM128 struct{}

// NROM256 is iNES mapper 000 with 2 PRG ROM banks
type NROM256 struct{}

// Translate performs mirroring and mapping of the address where needed
// and returns effective address
func (m *NROM128) Translate(addr uint16) uint16 {
	switch {
	case addr > 0xBFFF && addr <= 0xFFFF:
		return addr & 0xBFFF
	// PRG RAM mirroring $6000-$7FFF
	case addr > 0x6007 && addr <= 0x7FFF:
		return addr & 0x6007
	}
	return addr
}

// Translate performs mirroring and mapping of the address where needed
// and returns effective address
func (m *NROM256) Translate(addr uint16) uint16 {
	switch {
	// PRG RAM mirroring $6000-$7FFF
	case addr > 0x6007 && addr <= 0x7FFF:
		return addr & 0x6007
	}
	return addr
}
