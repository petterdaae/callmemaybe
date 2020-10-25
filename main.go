package main

import (
	"github.com/alecthomas/kong"
	"os"
	"os/exec"
)

type Arguments struct {
	Build Build `cmd:"build"`
}

type Build struct {
	File string `arg:"" type:"path"`
}

func (build *Build) Run() error {
	println("🔨 Building your executable with lang 1.0\n")

	nasmTemp := "out.nasm"
	oTemp := "out.o"

	content, err := ReadFile(build.File)
	if err != nil {
		return err
	}

	nasm, err := Compile(content)
	if err != nil {
		return nil
	}


	err = WriteFile(nasmTemp, nasm)
	if err != nil {
		return err
	}

	println("    Assembling ...")
	_, err = exec.Command("nasm", "-f", "elf64", nasmTemp,).CombinedOutput()
	if err != nil {
		println("❌ Assembling failed")
		return nil
	}

	println("    Linking ...")
	_, err = exec.Command("gcc", "-no-pie", "-o", "out", oTemp, "-lc").CombinedOutput()
	if err != nil {
		println("❌ Linking failed")
		return nil
	}

	os.Remove(nasmTemp)
	os.Remove(oTemp)

	println()
	println("🏆 Successful")

	return nil
}

func main() {
	var arguments Arguments
	ctx := kong.Parse(&arguments)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
