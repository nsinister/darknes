package ui

import (
	"darknes/common"
	"fmt"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowTitle = "DarkNES"

	fontPath = "assets/fonts/Poppins-Regular.ttf"
	fontSize = 18
)

type SdlFrontend struct {
	window  *sdl.Window
	surface *sdl.Surface

	cpuEmu common.CpuEmulator

	running bool

	// text overlay surface
	text *sdl.Surface
	font *ttf.Font
}

func CreateFrontend(cpuEmu common.CpuEmulator) *SdlFrontend {
	return &SdlFrontend{cpuEmu: cpuEmu}
}

func (frontend *SdlFrontend) renderText(textstr string, x int32, y int32) (err error) {
	if frontend.text, err = frontend.font.RenderUTF8Blended(textstr, sdl.Color{R: 255, G: 255, B: 255, A: 255}); err != nil {
		return err
	}
	defer frontend.text.Free()

	if err = frontend.text.Blit(nil, frontend.surface, &sdl.Rect{X: x, Y: y, W: 0, H: 0}); err != nil {
		return err
	}

	frontend.window.UpdateSurface()

	return nil
}

func (frontend *SdlFrontend) handleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		println("Quitting..")
		frontend.running = false
		break
	case *sdl.KeyboardEvent:
		if t.State == sdl.RELEASED {
			if t.Keysym.Sym == sdl.K_LEFT {
				// TODO:
			} else if t.Keysym.Sym == sdl.K_RIGHT {
				// TODO:
			}
			if t.Keysym.Sym == sdl.K_UP {
				// TODO:
			} else if t.Keysym.Sym == sdl.K_DOWN {
				// TODO:
			} else if t.Keysym.Sym == sdl.K_SPACE {
				// Perform one step in the CPU
				frontend.cpuEmu.Step()

				// Extract debug info
				cpu := frontend.cpuEmu.GetCpu()
				opName := strings.Split(cpu.LastOp.GetOpHandlerName(cpu.LastOp.Handler), ".")[1]

				cpuDebugText := fmt.Sprintf("OP: %s PC: [%04x] Status: [%08b] S: [(%x)] A: [%d (%x)] X: [%d (%x)] Y: [%d (%x)]",
					opName, cpu.PC, cpu.P, cpu.S, cpu.A, cpu.A, cpu.X, cpu.X, cpu.Y, cpu.Y)

				// Show some debug text
				frontend.surface.FillRect(nil, 0)
				frontend.renderText(cpuDebugText, 10, 10)
			}
		}
		break
	}
}

func (frontend *SdlFrontend) RunSdlLoop() (err error) {
	if err = ttf.Init(); err != nil {
		return
	}
	defer ttf.Quit()

	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		return
	}
	defer sdl.Quit()

	if frontend.window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN); err != nil {
		return
	}
	defer frontend.window.Destroy()

	if frontend.surface, err = frontend.window.GetSurface(); err != nil {
		return
	}

	// Load font for text output
	if frontend.font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		return
	}
	defer frontend.font.Close()

	frontend.window.UpdateSurface()

	// Main loop
	frontend.running = true
	for frontend.running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			frontend.handleEvent(event)
		}

		sdl.Delay(16)
	}

	return
}
