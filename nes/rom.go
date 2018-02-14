// Package rom implements iNES rom format
package nes

const (
	romHeaderLen uint32 = 16
)

// RomHeader represents header of the ROM file
type RomHeader struct {
	headerData []byte
	PrgRomSize byte
	ChrRomSize byte
	PrgRAMSize byte
	Flags6     byte
	Flags7     byte
	Flags9     byte
	HasTrainer bool
	MapperNum  byte
	PrgBegin   uint32
	PrgEnd     uint32
}

// Rom contains ordered ROM data fields in iNES format
type Rom struct {
	data   []byte
	Header *RomHeader
	prgRom []byte
	chrRom []byte
}

// Read method returns byte from ROM at the specified address
func (r *Rom) Read(addr uint32) byte {
	return r.data[addr]
}

// LoadRomData parses raw ROM data and transforms into NesRom struct
func LoadRomData(romData []byte) *Rom {
	var rom = Rom{
		data: romData,
		Header: &RomHeader{
			headerData: romData[:romHeaderLen],
			PrgRomSize: romData[4],
			ChrRomSize: romData[5],
			PrgRAMSize: romData[8],
			Flags6:     romData[6],
			Flags7:     romData[7],
			Flags9:     romData[9],
			HasTrainer: (romData[6] & 4) != 0,
			MapperNum:  ((romData[6] & 0xf0) >> 4) | (romData[7] & 0xf0),
			PrgBegin:   romHeaderLen,
		},
	}

	if rom.Header.HasTrainer {
		rom.Header.PrgBegin += 512
	}

	rom.Header.PrgEnd = rom.Header.PrgBegin + (uint32(rom.Header.PrgRomSize) * (16 * 1024))

	rom.prgRom = rom.data[rom.Header.PrgBegin:rom.Header.PrgEnd]

	if rom.Header.ChrRomSize > 0 {
		chrEnd := rom.Header.PrgEnd + (uint32(rom.Header.ChrRomSize) * 8192)
		rom.chrRom = rom.data[rom.Header.PrgEnd:chrEnd]
	}
	return &rom
}
