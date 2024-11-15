package nes

const (
	// PPU register addresses
	PPUController uint16 = 0x2000
	PPUMask       uint16 = 0x2001
	PPUStatus     uint16 = 0x2002
	OAMAddress    uint16 = 0x2003
	OAMData       uint16 = 0x2004
	PPUScroll     uint16 = 0x2005
	PPUAddress    uint16 = 0x2006
	PPUData       uint16 = 0x2007
	OAMDMA        uint16 = 0x4014

	// PPU Status register flags
	PPUStatusVBlank byte = 0x80
)

type PPU struct {
	cpu      *CPU
	cycles   uint16
	scanline uint16
	nmi      bool
}

func InitPPU(c *CPU) *PPU {
	return &PPU{cpu: c}
}

func (ppu *PPU) Reset() {
	ppu.cycles = 0
	ppu.nmi = false
	ppu.scanline = 0
}

func (ppu *PPU) Step(cycles uint16) {
	ppu.cycles += cycles
	if ppu.cycles >= 341 {
		ppu.cycles -= 341
		ppu.scanline++
	}

	if 0 <= ppu.scanline && ppu.scanline <= 239 {

	} else if ppu.scanline == 241 && ppu.cycles == 1 {
		// VBlank
		println("VBLANK!!!")
		ppu.setStatus(PPUStatusVBlank)
		ppu.nmi = true

		// FIXME:
		ppu.cpu.nmiInterrupt()
		// TODO

	} else if ppu.scanline == 261 && ppu.cycles == 1 {
		// VBlank off
		println("vblank off")
		ppu.clearStatus(PPUStatusVBlank)
		ppu.nmi = false
		ppu.scanline = 0
	}

}

func (ppu *PPU) setStatus(flag byte) {
	status := ppu.cpu.mem.Read(PPUStatus)
	ppu.cpu.mem.Write(PPUStatus, status|flag)
}

func (ppu *PPU) clearStatus(flag byte) {
	status := ppu.cpu.mem.Read(PPUStatus)
	ppu.cpu.mem.Write(PPUStatus, status & ^flag)
}
