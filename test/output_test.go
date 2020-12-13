package test

import (
	"callmemaybe/utils"
	"testing"
)

func TestCase001(t *testing.T) {
	utils.AssertProgramOutput("testcases/001.cmm", "1\n", t)
}

func TestCase002(t *testing.T) {
	utils.AssertProgramOutput("testcases/002.cmm", "7\n", t)
}

func TestCase003(t *testing.T) {
	utils.AssertProgramOutput("testcases/003.cmm", "15\n", t)
}

func TestCase004(t *testing.T) {
	utils.AssertProgramOutput("testcases/004.cmm", "12\n", t)
}

func TestCase005(t *testing.T) {
	utils.AssertProgramOutput("testcases/005.cmm", "", t)
}

func TestCase006(t *testing.T) {
	utils.AssertProgramOutput("testcases/006.cmm", "", t)
}

func TestCase007(t *testing.T) {
	utils.AssertCompilerFails("testcases/007.cmm", t)
}

func TestCase008(t *testing.T) {
	utils.AssertCompilerFails("testcases/008.cmm", t)
}

func TestCase009(t *testing.T) {
	utils.AssertCompilerFails("testcases/009.cmm", t)
}

func TestCase010(t *testing.T) {
	utils.AssertCompilerFails("testcases/010.cmm", t)
}

func TestCase011(t *testing.T) {
	utils.AssertProgramOutput("testcases/011.cmm", "5\n", t)
}

func TestCase012(t *testing.T) {
	utils.AssertCompilerFails("testcases/012.cmm", t)
}

func TestCase013(t *testing.T) {
	utils.AssertProgramOutput("testcases/013.cmm", "3\n2\n", t)
}

func TestCase014(t *testing.T) {
	utils.AssertProgramOutput("testcases/014.cmm", "3\n2\n1\n0\n", t)
}

func TestCase015(t *testing.T) {
	utils.AssertProgramOutput("testcases/015.cmm", "42\n", t)
}

func TestCase016(t *testing.T) {
	utils.AssertProgramOutput("testcases/016.cmm", "42\n42\n", t)
}

func TestCase017(t *testing.T) {
	utils.AssertProgramOutput("testcases/017.cmm", "50\n", t)
}

func TestCase018(t *testing.T) {
	utils.AssertProgramOutput("testcases/018.cmm", "48\n", t)
}

func TestCase019(t *testing.T) {
	utils.AssertProgramOutput("testcases/019.cmm", "5\n", t)
}

func TestCase020(t *testing.T) {
	utils.AssertProgramOutput("testcases/020.cmm", "21\n", t)
}

func TestCase021(t *testing.T) {
	utils.AssertProgramOutput("testcases/021.cmm", "2\n", t)
}

func TestCase022(t *testing.T) {
	utils.AssertProgramOutput("testcases/022.cmm", "12\n", t)
}

func TestCase023(t *testing.T) {
	utils.AssertProgramOutput("testcases/023.cmm", "230\n", t)
}

func TestCase024(t *testing.T) {
	utils.AssertCompilerFails("testcases/024.cmm", t)
}

func TestCase025(t *testing.T) {
	utils.AssertCompilerFails("testcases/025.cmm", t)
}

func TestCase026(t *testing.T) {
	utils.AssertProgramOutput("testcases/026.cmm", "10\n", t)
}

func TestCase027(t *testing.T) {
	utils.AssertProgramOutput("testcases/027.cmm", "6\n", t)
}

func TestCase028(t *testing.T) {
	utils.AssertProgramOutput("testcases/028.cmm", "32\n", t)
}

func TestCase029(t *testing.T) {
	utils.AssertProgramOutput("testcases/029.cmm", "16\n", t)
}

func TestCase030(t *testing.T) {
	utils.AssertProgramOutput("testcases/030.cmm", "7\n", t)
}

func TestCase031(t *testing.T) {
	utils.AssertProgramOutput("testcases/031.cmm", "20\n", t)
}

func TestCase032(t *testing.T) {
	utils.AssertProgramOutput("testcases/032.cmm", "1\n", t)
}

func TestCase033(t *testing.T) {
	utils.AssertProgramOutput("testcases/033.cmm", "0\n", t)
}

func TestCase034(t *testing.T) {
	utils.AssertProgramOutput("testcases/034.cmm", "1\n", t)
}

func TestCase035(t *testing.T) {
	utils.AssertProgramOutput("testcases/035.cmm", "1\n", t)
}

func TestCase036(t *testing.T) {
	utils.AssertCompilerFails("testcases/036.cmm", t)
}

func TestCase037(t *testing.T) {
	utils.AssertProgramOutput("testcases/037.cmm", "1\n0\n0\n0\n", t)
}

func TestCase038(t *testing.T) {
	utils.AssertProgramOutput("testcases/038.cmm", "3\n", t)
}

func TestCase039(t *testing.T) {
	utils.AssertProgramOutput("testcases/039.cmm", "0\n1\n", t)
}

func TestCase040(t *testing.T) {
	utils.AssertProgramOutput("testcases/040.cmm", "-5\n-4\n-3\n", t)
}

func TestCase041(t *testing.T) {
	utils.AssertProgramOutput("testcases/041.cmm", "3\n-3\n", t)
}

func TestCase042(t *testing.T) {
	utils.AssertProgramOutput("testcases/042.cmm", "24\n", t)
}

func TestCase043(t *testing.T) {
	utils.AssertProgramOutput("testcases/043.cmm", "0\n2\n1\n", t)
}

func TestCase044(t *testing.T) {
	utils.AssertProgramCrashes("testcases/044.cmm", t)
}

func TestCase045(t *testing.T) {
	utils.AssertProgramOutput("testcases/045.cmm", "-1\n-2\n", t)
}

func TestCase046(t *testing.T) {
	utils.AssertProgramOutput("testcases/046.cmm", "1\n", t)
}

func TestCase047(t *testing.T) {
	utils.AssertProgramOutput("testcases/047.cmm", "5\n4\n3\n2\n1\n42\n", t)
}

func TestCase048(t *testing.T) {
	utils.AssertProgramOutput("testcases/048.cmm", "1\n1\n0\n0\n1\n", t)
}

func TestCase049(t *testing.T) {
	utils.AssertProgramOutput("testcases/049.cmm", "a\n", t)
}

func TestCase050(t *testing.T) {
	utils.AssertProgramOutput("testcases/050.cmm", "H\ne\nl\nl\no\n \nw\no\nr\nl\nd\n!\n", t)
}

func TestCase051(t *testing.T) {
	utils.AssertProgramOutput("testcases/051.cmm", "p\np\np\np\np\np\np\np\np\np\n", t)
}

func TestCase052(t *testing.T) {
	utils.AssertProgramOutput("testcases/052.cmm", "2\n", t)
}

func TestCase053(t *testing.T) {
	utils.AssertProgramOutput("testcases/053.cmm", "o\nd\n", t)
}

func TestCase054(t *testing.T) {
	utils.AssertCompilerFails("testcases/054.cmm", t)
}

func TestCase055(t *testing.T) {
	utils.AssertCompilerFails("testcases/055.cmm", t)
}

func TestCase056(t *testing.T) {
	utils.AssertProgramOutput("testcases/056.cmm", "3\n", t)
}

func TestCase057(t *testing.T) {
	utils.AssertProgramOutput("testcases/057.cmm", "l\n", t)
}

func TestCase058(t *testing.T) {
	utils.AssertProgramOutput("testcases/058.cmm", "\\\n", t)
}
