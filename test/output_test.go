package test

import (
	"lang/utils"
	"testing"
)

func TestCase001(t *testing.T) {
	utils.AssertProgramOutput("testcases/001.lang", "1\n", t)
}

func TestCase002(t *testing.T) {
	utils.AssertProgramOutput("testcases/002.lang", "7\n", t)
}

func TestCase003(t *testing.T) {
	utils.AssertProgramOutput("testcases/003.lang", "15\n", t)
}

func TestCase004(t *testing.T) {
	utils.AssertProgramOutput("testcases/004.lang", "12\n", t)
}
