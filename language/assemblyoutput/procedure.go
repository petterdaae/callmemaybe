package assemblyoutput

import (
	"callmemaybe/language/common"
)

type procedure struct {
	Name                              string
	NumberOfArgs                      int
	StackSizeBeforeFunctionGeneration int
	Operations                        []string
	ReturnKind                        common.ContextElementKind
	Args                              []common.Arg
}
