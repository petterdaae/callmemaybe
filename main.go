package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"os"
	"os/exec"
	"callmemaybe/utils"
)

type Arguments struct {
	Build Build `cmd:"build"`
	Nasm Nasm `cmd:"nasm"`
}

type Build struct {
	File string `arg:"" type:"path"`
}

type Nasm struct {
	File string `arg:"" type:"path"`
}

func (build *Build) Run() error {
	println("🔨 Building your executable with callmemaybe 1.0\n")

	nasmTemp := "out.nasm"
	oTemp := "out.o"

	content, err := utils.ReadFile(build.File)
	if err != nil {
		return err
	}

	println("    Compiling ...")
	nasm, err := utils.Compile(content)
	if err != nil {
		println("❌ Failed to compile")
		println(err.Error())
		return nil
	}

	err = utils.WriteFile(nasmTemp, nasm)
	if err != nil {
		return err
	}

	println("    Assembling ...")
	_, err = exec.Command("nasm", "-f", "elf64", nasmTemp).CombinedOutput()
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

	println("\n🏆 Successful")

	return nil
}

func (args *Nasm) Run() error {
	content, err := utils.ReadFile(args.File)
	if err != nil {
		return err
	}
	nasm, err := utils.Compile(content)
	if err != nil {
		return nil
	}
	fmt.Println(nasm)
	return nil
}

func main() {
	var arguments Arguments
	ctx := kong.Parse(&arguments)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
