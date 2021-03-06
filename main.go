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
	X86   X86   `cmd:"x86"`
}

type Build struct {
	File string `arg:"" type:"path"`
}

type X86 struct {
	File string `arg:"" type:"path"`
}

func (build *Build) Run() error {
	nasmTemp := "out.nasm"
	oTemp := "out.o"
	content, err := utils.ReadFile(build.File)
	if err != nil {
		return err
	}
	nasm, err := utils.Compile(content)
	if err != nil {
		println(err.Error())
		return nil
	}
	err = utils.WriteFile(nasmTemp, nasm)
	if err != nil {
		return err
	}
	_, err = exec.Command("nasm", "-f", "elf64", nasmTemp).CombinedOutput()
	if err != nil {
		println(err.Error())
		return nil
	}
	_, err = exec.Command("gcc", "-no-pie", "-o", "out", oTemp, "-lc").CombinedOutput()
	if err != nil {
		println(err.Error())
		return nil
	}
	os.Remove(nasmTemp)
	os.Remove(oTemp)
	return nil
}

func (args *X86) Run() error {
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
