package main

import (
	"fmt"
	"os"

	"darknes/nes"
	"darknes/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No rom file specified")
		return
	}
	path := os.Args[1]
	romData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	r := nes.LoadRomData(romData)
	mem := r.Load()
	cpu := nes.InitCPU(mem)
	cpu.Reset()

	fmt.Printf("Init state: A=%x, X=%x, Y=%x, S=%x, P=%b, PC=%x\n",
		cpu.A, cpu.X, cpu.Y, cpu.S, cpu.P, cpu.PC)

	screen := ui.CreateScreen(cpu)
	if err := screen.RunSdlLoop(); err != nil {
		os.Exit(1)
	}

	/*
		for i := 0; i < 10; i++ {
			cpu.Step()
			fmt.Printf("A=%x, X=%x, Y=%x, S=%x, P=%b, PC=%x\n",
				cpu.A, cpu.X, cpu.Y, cpu.S, cpu.P, cpu.PC)
		}
	*/
}
