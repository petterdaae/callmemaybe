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

func TestCase005(t *testing.T) {
	utils.AssertProgramOutput("testcases/005.lang", "", t)
}

func TestCase006(t *testing.T) {
	utils.AssertProgramOutput("testcases/006.lang", "", t)
}

func TestCase007(t *testing.T) {
	utils.AssertProgramOutput("testcases/007.lang", "43\n", t)
}

func TestCase008(t *testing.T) {
	utils.AssertProgramOutput("testcases/008.lang", "43\n", t)
}

func TestCase009(t *testing.T) {
	utils.AssertProgramOutput("testcases/009.lang", "43\n43\n45\n45\n", t)
}

func TestCase010(t *testing.T) {
	utils.AssertProgramOutput("testcases/010.lang", "43\n", t)
}

func TestCase011(t *testing.T) {
	utils.AssertProgramOutput("testcases/011.lang", "5\n", t)
}

func TestCase012(t *testing.T) {
	utils.AssertProgramOutput("testcases/012.lang", "5\n", t)
}

func TestCase013(t *testing.T) {
	utils.AssertProgramOutput("testcases/013.lang", "3\n2\n", t)
}

func TestCase014(t *testing.T) {
	utils.AssertProgramOutput("testcases/014.lang", "3\n2\n1\n0\n", t)
}

func TestCase015(t *testing.T) {
	utils.AssertProgramOutput("testcases/015.lang", "42\n", t)
}

func TestCase016(t *testing.T) {
	utils.AssertProgramOutput("testcases/016.lang", "42\n42\n", t)
}

func TestCase017(t *testing.T) {
	utils.AssertProgramOutput("testcases/017.lang", "50\n", t)
}

func TestCase018(t *testing.T) {
	utils.AssertProgramOutput("testcases/018.lang", "48\n", t)
}

func TestCase019(t *testing.T) {
	utils.AssertProgramOutput("testcases/019.lang", "5\n", t)
}

func TestCase020(t *testing.T) {
	utils.AssertProgramOutput("testcases/020.lang", "21\n", t)
}

func TestCase021(t *testing.T) {
	utils.AssertProgramOutput("testcases/021.lang", "2\n", t)
}

func TestCase022(t *testing.T) {
	utils.AssertProgramOutput("testcases/022.lang", "12\n", t)
}

func TestCase023(t *testing.T) {
	utils.AssertProgramOutput("testcases/023.lang", "230\n", t)
}

func TestCase024(t *testing.T) {
	utils.AssertProgramOutput("testcases/024.lang", "100\n", t)
}

func TestCase025(t *testing.T) {
	utils.AssertCompilerFails("testcases/025.lang", t)
}
