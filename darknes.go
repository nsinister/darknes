package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/skyphaser/darknes/nes"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No rom file specified")
		return
	}
	path := os.Args[1]
	romData, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	r := nes.LoadRomData(romData)
	mem := r.Load()
	cpu := nes.InitCPU(mem)
	cpu.Reset()
	fmt.Printf("Init state: A=%x, X=%x, Y=%x, S=%x, P=%b, PC=%x\n",
		cpu.A, cpu.X, cpu.Y, cpu.S, cpu.P, cpu.PC)
	fmt.Printf("%d\n", r.Header.MapperNum)
	fmt.Printf("%x\n", mem.Translate(0x2009))
	for i := 0; i < 4; i++ {
		cpu.Step()
		fmt.Printf("A=%x, X=%x, Y=%x, S=%x, P=%b, PC=%x\n",
			cpu.A, cpu.X, cpu.Y, cpu.S, cpu.P, cpu.PC)
	}
}
